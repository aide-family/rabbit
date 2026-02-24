package impl

import (
	"context"
	"strings"
	"sync"
	"time"

	magicboxapiv1 "github.com/aide-family/magicbox/api/v1"
	"github.com/aide-family/magicbox/config"
	"github.com/aide-family/magicbox/connect"
	"github.com/aide-family/magicbox/contextx"
	"github.com/aide-family/magicbox/enum"
	"github.com/aide-family/magicbox/merr"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/query"
	"github.com/aide-family/rabbit/internal/data/impl/state"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/message"
)

func NewMessageRepository(
	c *conf.Bootstrap,
	d *data.Data,
	messageLogRepo repository.MessageLog,
	namespaceRepo repository.Namespace,
) (repository.Message, error) {
	jobCore := c.GetJobCore()
	repo := &messageRepository{
		Data:           d,
		messageLogRepo: messageLogRepo,
		namespaceRepo:  namespaceRepo,
		messageChan:    make(chan *state.MessageTask, jobCore.GetBufferSize()),
		stopChan:       make(chan struct{}),
		workerTotal:    int(jobCore.GetWorkerTotal()),
		timeout:        jobCore.GetTimeout().AsDuration(),
		wg:             sync.WaitGroup{},
	}
	state.RegisterMessageTaskProcess(enum.MessageStatus_MessageStatus_UNKNOWN, repo.unknownMessageTaskProcess)
	state.RegisterMessageTaskProcess(enum.MessageStatus_PENDING, repo.pendingMessageTaskProcess)
	state.RegisterMessageTaskProcess(enum.MessageStatus_SENDING, repo.sendingMessageTaskProcess)
	state.RegisterMessageTaskProcess(enum.MessageStatus_SENT, repo.sentMessageTaskProcess)
	state.RegisterMessageTaskProcess(enum.MessageStatus_FAILED, repo.failedMessageTaskProcess)
	state.RegisterMessageTaskProcess(enum.MessageStatus_CANCELLED, repo.cancelledMessageTaskProcess)
	for _, value := range enum.MessageStatus_value {
		status := enum.MessageStatus(value)
		state.RegisterMessageTaskState(status, state.NewMessageTaskState(status))
	}

	query.SetDefault(d.DB())
	if err := repo.initClusters(c.GetJobClusters()); err != nil {
		return nil, err
	}
	if err := repo.Start(context.Background()); err != nil {
		return nil, err
	}
	if err := repo.loadMessageLogs(); err != nil {
		return nil, err
	}
	d.AppendClose("messageRepo", func() error { return repo.Stop(context.Background()) })
	return repo, nil
}

type messageRepository struct {
	messageLogRepo repository.MessageLog
	namespaceRepo  repository.Namespace
	stopChan       chan struct{}
	messageChan    chan *state.MessageTask
	wg             sync.WaitGroup
	workerTotal    int
	timeout        time.Duration
	clustersMu     sync.RWMutex
	clusters       []ClusterSender
	*data.Data
}

func (m *messageRepository) pendingMessageTaskProcess(task *state.MessageTask) (enum.MessageStatus, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	ctx = contextx.WithNamespace(ctx, task.NamespaceUID)
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, task.MessageUID)
	if err != nil {
		return 0, false
	}

	if messageLog.Status != enum.MessageStatus_PENDING {
		return 0, false
	}
	changed, err := m.messageLogRepo.UpdateMessageLogStatusSendingIf(ctx, messageLog.UID, enum.MessageStatus_PENDING)
	if err != nil {
		return 0, false
	}

	return enum.MessageStatus_SENDING, changed
}

func (m *messageRepository) sendingMessageTaskProcess(task *state.MessageTask) (enum.MessageStatus, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	ctx = contextx.WithNamespace(ctx, task.NamespaceUID)
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, task.MessageUID)
	if err != nil {
		return 0, false
	}

	if messageLog.Status != enum.MessageStatus_SENDING {
		return 0, false
	}

	driver, ok := message.GetDriver(messageLog.MessageType)
	if !ok {
		return 0, false
	}
	msgConfig, err := messageLog.ToMessageConfig()
	if err != nil {
		return 0, false
	}
	sender, err := driver(msgConfig)
	if err != nil {
		return 0, false
	}
	msg := message.NewMessage(messageLog.MessageType, []byte(messageLog.Message))
	if err := sender.Send(ctx, msg); err != nil {
		changed, err := m.messageLogRepo.UpdateMessageLogLastErrorIf(ctx, task.MessageUID, enum.MessageStatus_SENDING, err.Error())
		if err != nil {
			return 0, false
		}
		return enum.MessageStatus_FAILED, changed
	}
	changed, err := m.messageLogRepo.UpdateMessageLogStatusSuccessIf(ctx, task.MessageUID)
	if err != nil {
		return 0, false
	}
	return enum.MessageStatus_SENT, changed
}

func (m *messageRepository) sentMessageTaskProcess(task *state.MessageTask) (enum.MessageStatus, bool) {
	return 0, false
}

func (m *messageRepository) cancelledMessageTaskProcess(task *state.MessageTask) (enum.MessageStatus, bool) {
	return 0, false
}

func (m *messageRepository) failedMessageTaskProcess(task *state.MessageTask) (enum.MessageStatus, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	ctx = contextx.WithNamespace(ctx, task.NamespaceUID)
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, task.MessageUID)
	if err != nil {
		return 0, false
	}
	if messageLog.Status != enum.MessageStatus_FAILED || task.IsMaxRetry() {
		return 0, false
	}
	changed, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, task.MessageUID, enum.MessageStatus_FAILED, enum.MessageStatus_PENDING)
	if err != nil {
		return 0, false
	}

	return enum.MessageStatus_PENDING, changed
}

func (m *messageRepository) unknownMessageTaskProcess(task *state.MessageTask) (enum.MessageStatus, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	ctx = contextx.WithNamespace(ctx, task.NamespaceUID)
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, task.MessageUID)
	if err != nil {
		return 0, false
	}
	if messageLog.Status == enum.MessageStatus_MessageStatus_UNKNOWN {
		changed, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, task.MessageUID, enum.MessageStatus_MessageStatus_UNKNOWN, enum.MessageStatus_PENDING)
		if err != nil {
			return 0, false
		}
		return enum.MessageStatus_PENDING, changed
	}

	return messageLog.Status, true
}

// AppendMessage implements [repository.Message].
func (m *messageRepository) AppendMessage(ctx context.Context, messageUID snowflake.ID) error {
	task := &state.MessageTask{
		NamespaceUID: contextx.GetNamespace(ctx),
		MessageUID:   messageUID,
	}
	select {
	case m.messageChan <- task:
		klog.Context(ctx).Debugw("msg", "append message success", "messageUID", messageUID)
		return nil
	default:
		klog.Context(ctx).Debugw("msg", "append message channel full", "messageUID", messageUID)
		if task.IsMaxRetry() {
			klog.Context(ctx).Warnw("msg", "append message retry count reached max", "task", task)
			return nil
		}
		m.clustersMu.RLock()
		defer m.clustersMu.RUnlock()
		for _, cluster := range m.clusters {
			if err := cluster.Send(ctx, messageUID); err != nil {
				klog.Warnw("msg", "send message to cluster failed", "error", err, "cluster", cluster)
			}
		}
		klog.Context(ctx).Debugw("msg", "append message retry", "task", task)
		task.RetryIncrement()
		return m.AppendMessage(ctx, messageUID)
	}
}

func (m *messageRepository) loadMessageLogs() error {
	req := &magicboxapiv1.SelectNamespaceRequest{
		Status: enum.GlobalStatus_ENABLED,
		Limit:  1000,
	}
	wg := sync.WaitGroup{}
	for {
		selectNamespaceBoResult, err := m.namespaceRepo.SelectNamespace(context.Background(), req)
		if err != nil {
			return err
		}

		wg.Go(func() {
			bufferChan := make(chan struct{}, 3)
			for _, namespace := range selectNamespaceBoResult.Items {
				bufferChan <- struct{}{}
				go func(namespaceUID snowflake.ID) {
					m.loadMessageLogsForNamespace(namespaceUID)
					<-bufferChan
				}(snowflake.ParseInt64(namespace.Value))
			}
		})
		if !selectNamespaceBoResult.HasMore {
			break
		}
		req.LastUID = selectNamespaceBoResult.LastUID
	}

	wg.Wait()
	return nil
}

func (m *messageRepository) loadMessageLogsForNamespace(namespaceUID snowflake.ID) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = contextx.WithNamespace(ctx, namespaceUID)
	messageLogs, err := m.messageLogRepo.GetAllMessageLogs(ctx, enum.MessageStatus_PENDING)
	if err != nil {
		klog.Warnw("msg", "load message logs for namespace failed", "error", err, "namespaceUID", namespaceUID)
		return
	}
	for _, messageLog := range messageLogs {
		m.messageChan <- &state.MessageTask{
			NamespaceUID: messageLog.NamespaceUID,
			MessageUID:   messageLog.UID,
		}
	}
}

func (m *messageRepository) worker(number int) {
	for {
		select {
		case <-m.stopChan:
			klog.Debugw("msg", "message worker stopped", "workerIndex", number)
			return
		case task := <-m.messageChan:
			m.processMessageTask(task)
		}
	}
}

func (m *messageRepository) processMessageTask(task *state.MessageTask) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	ctx = contextx.WithNamespace(ctx, task.NamespaceUID)
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, task.MessageUID)
	if err != nil {
		klog.Warnw("msg", "get message log failed", "error", err, "messageUID", task.MessageUID)
		return
	}
	messageTaskState, ok := state.GetMessageTaskState(messageLog.Status)
	if !ok {
		klog.Warnw("msg", "message task state not found", "status", messageLog.Status)
		return
	}
	messageTaskState.Process(ctx, task)
}

// Start implements [repository.Message].
func (m *messageRepository) Start(ctx context.Context) error {
	klog.Infow("msg", "start message worker", "workerTotal", m.workerTotal)
	for i := 0; i < m.workerTotal; i++ {
		index := i
		m.wg.Go(func() {
			m.worker(index)
		})
	}
	return nil
}

// Stop implements [repository.Message].
func (m *messageRepository) Stop(_ context.Context) error {
	close(m.stopChan)
	m.wg.Wait()
	klog.Infow("msg", "message worker stopped")
	m.clustersMu.RLock()
	defer m.clustersMu.RUnlock()
	for _, cluster := range m.clusters {
		if err := cluster.Close(); err != nil {
			klog.Warnw("msg", "close cluster failed", "error", err, "name", cluster.Name())
		}
		klog.Debugw("msg", "close cluster success", "name", cluster.Name())
	}
	klog.Infow("msg", "close clusters success")
	return nil
}

type ClusterSender interface {
	Name() string
	Send(ctx context.Context, messageUID snowflake.ID) error
	Close() error
}

func NewClusterSender(name string, sendFunc func(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error), closeFunc func() error) ClusterSender {
	return &clusterSender{
		name:      name,
		sendFunc:  sendFunc,
		closeFunc: closeFunc,
	}
}

type clusterSender struct {
	name      string
	sendFunc  func(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error)
	closeFunc func() error
}

func (c *clusterSender) Name() string {
	return c.name
}

func (c *clusterSender) Send(ctx context.Context, messageUID snowflake.ID) error {
	req := &apiv1.SendMessageRequest{
		Uid: messageUID.Int64(),
	}
	_, err := c.sendFunc(ctx, req)
	return err
}

func (c *clusterSender) Close() error {
	return c.closeFunc()
}

func (m *messageRepository) initClusters(c *config.ClusterConfig) error {
	clusterEndpoints := strutil.SplitSkipEmpty(c.GetEndpoints(), ",")
	clusterTimeout := c.GetTimeout().AsDuration()
	clusterName := c.GetName()
	protocol := c.GetProtocol().String()
	clusters := make([]ClusterSender, 0, len(clusterEndpoints))
	for _, clusterEndpoint := range clusterEndpoints {
		opts := []connect.InitOption{
			connect.WithDiscovery(m.Registry()),
		}
		initConfig := connect.NewDefaultConfig(clusterName, clusterEndpoint, clusterTimeout, protocol)

		name := strings.Join([]string{clusterName, clusterEndpoint}, ":")
		var clusterSender ClusterSender
		switch protocol {
		case connect.ProtocolHTTP:
			httpClient, err := connect.InitHTTPClient(initConfig, opts...)
			if err != nil {
				klog.Warnw("msg", "create HTTP client failed", "endpoint", clusterEndpoint, "error", err)
				return err
			}

			httpSender := apiv1.NewSenderHTTPClient(httpClient)
			clusterSender = NewClusterSender(name, func(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error) {
				return httpSender.SendMessage(ctx, req)
			}, httpClient.Close)
		case connect.ProtocolGRPC:
			grpcClient, err := connect.InitGRPCClient(initConfig, opts...)
			if err != nil {
				klog.Warnw("msg", "create GRPC client failed", "endpoint", clusterEndpoint, "error", err)
				return err
			}
			grpcSender := apiv1.NewSenderClient(grpcClient)
			clusterSender = NewClusterSender(name, func(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error) {
				return grpcSender.SendMessage(ctx, req)
			}, grpcClient.Close)
		default:
			klog.Warnw("msg", "unknown protocol", "endpoint", clusterEndpoint, "protocol", protocol)
			return merr.ErrorInternalServer("unknown protocol: %s", protocol)
		}

		clusters = append(clusters, clusterSender)
	}
	m.clustersMu.Lock()
	defer m.clustersMu.Unlock()
	m.clusters = clusters
	return nil
}
