package impl

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/aide-family/magicbox/config"
	"github.com/aide-family/magicbox/connect"
	"github.com/aide-family/magicbox/contextx"
	"github.com/aide-family/magicbox/enum"
	"github.com/aide-family/magicbox/merr"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/query"
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
		messageChan:    make(chan *messageTask, jobCore.GetBufferSize()),
		stopChan:       make(chan struct{}),
		workerTotal:    int(jobCore.GetWorkerTotal()),
		timeout:        jobCore.GetTimeout().AsDuration(),
		wg:             sync.WaitGroup{},
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
	messageChan    chan *messageTask
	wg             sync.WaitGroup
	workerTotal    int
	timeout        time.Duration
	clustersMu     sync.RWMutex
	clusters       []ClusterSender
	*data.Data
}

type messageTask struct {
	namespaceUID snowflake.ID
	messageUID   snowflake.ID
}

// AppendMessage implements [repository.Message].
func (m *messageRepository) AppendMessage(ctx context.Context, messageUID snowflake.ID) error {
	task := &messageTask{
		namespaceUID: contextx.GetNamespace(ctx),
		messageUID:   messageUID,
	}
	select {
	case m.messageChan <- task:
		klog.Context(ctx).Debugw("msg", "append message success", "messageUID", messageUID)
		return nil
	default:
		klog.Context(ctx).Warnw("msg", "append message channel full", "messageUID", messageUID)
		// TODO: send message to other services
		m.clustersMu.RLock()
		defer m.clustersMu.RUnlock()
		for _, cluster := range m.clusters {
			if err := cluster.Send(ctx, messageUID); err != nil {
				klog.Warnw("msg", "send message to cluster failed", "error", err, "cluster", cluster)
			}
		}
		return merr.ErrorInternalServer("append message channel full: %d", messageUID)
	}
}

// SendMessage implements [repository.Message].
func (m *messageRepository) SendMessage(ctx context.Context, messageUID snowflake.ID) error {
	return m.sendMessage(contextx.GetNamespace(ctx), messageUID)
}

// SendMessage implements [repository.Message].
func (m *messageRepository) sendMessage(namespaceUID, messageUID snowflake.ID) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	ctx = contextx.WithNamespace(ctx, namespaceUID)
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, messageUID)
	if err != nil {
		return err
	}
	if slices.Contains([]enum.MessageStatus{enum.MessageStatus_SENT, enum.MessageStatus_SENDING, enum.MessageStatus_CANCELLED}, messageLog.Status) {
		klog.Context(ctx).Warnw("msg", "message already sent or sending or cancelled", "messageUID", messageUID, "status", messageLog.Status)
		return nil
	}
	changed, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, messageUID, messageLog.Status, enum.MessageStatus_SENDING)
	if err != nil {
		return err
	}
	if !changed {
		klog.Context(ctx).Warnw("msg", "message status is not sending, message sending failed", "messageUID", messageUID)
		return merr.ErrorNotFound("message sending failed, the status of this message has changed")
	}

	return m.processMessage(ctx, messageLog)
}

func (m *messageRepository) loadMessageLogs() error {
	req := &bo.SelectNamespaceBo{
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
		m.messageChan <- &messageTask{
			namespaceUID: messageLog.NamespaceUID,
			messageUID:   messageLog.UID,
		}
	}
}

// messageLogBody 实现 message.Message，用于把 messageLog 的 Message 字段传给 sender.Send。
type messageLogBody struct {
	messageType enum.MessageType
	body        []byte
}

func (m *messageLogBody) Type() enum.MessageType   { return m.messageType }
func (m *messageLogBody) Marshal() ([]byte, error) { return m.body, nil }

func (m *messageRepository) processMessage(ctx context.Context, messageLog *bo.MessageLogItemBo) error {
	driver, ok := message.GetDriver(messageLog.MessageType)
	if !ok {
		klog.Context(ctx).Warnw("msg", "message driver not found", "messageType", messageLog.MessageType)
		return merr.ErrorNotFound("message driver not found: %d", messageLog.MessageType)
	}

	msgConfig, err := messageLog.ToMessageConfig()
	if err != nil {
		klog.Context(ctx).Warnw("msg", "message config convert failed", "error", err, "messageType", messageLog.MessageType)
		return err
	}

	sender, err := driver(msgConfig)
	if err != nil {
		klog.Context(ctx).Warnw("msg", "message driver create failed", "error", err, "messageType", messageLog.MessageType)
		return err
	}

	msg := &messageLogBody{
		messageType: messageLog.MessageType,
		body:        []byte(string(messageLog.Message)),
	}
	if err := sender.Send(ctx, msg); err != nil {
		klog.Context(ctx).Warnw("msg", "message send failed", "error", err, "messageType", messageLog.MessageType)
		changed, err := m.messageLogRepo.UpdateMessageLogLastErrorIf(ctx, messageLog.UID, enum.MessageStatus_SENDING, err.Error())
		if err != nil {
			klog.Context(ctx).Warnw("msg", "message log last error update failed", "error", err, "messageUID", messageLog.UID)
			return err
		}
		if !changed {
			klog.Context(ctx).Warnw("msg", "message last error is not sending, message sending failed", "messageUID", messageLog.UID)
			return merr.ErrorNotFound("message sending failed, the last error of this message has changed")
		}
		return err
	}
	changed, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, messageLog.UID, enum.MessageStatus_SENDING, enum.MessageStatus_SENT)
	if err != nil {
		klog.Context(ctx).Warnw("msg", "message status update failed", "error", err, "messageUID", messageLog.UID, "status", enum.MessageStatus_SENDING, "newStatus", enum.MessageStatus_SENT)
		return err
	}
	if !changed {
		klog.Context(ctx).Warnw("msg", "message status is not sending, message sending failed", "messageUID", messageLog.UID)
		return merr.ErrorNotFound("message sending failed, the status of this message has changed")
	}
	return nil
}

// Start implements [repository.Message].
func (m *messageRepository) Start(ctx context.Context) error {
	klog.Infow("msg", "start message worker", "workerTotal", m.workerTotal)
	for i := 0; i < m.workerTotal; i++ {
		index := i
		m.wg.Go(func() {
			for {
				select {
				case <-m.stopChan:
					klog.Debugw("msg", "message worker stopped", "workerIndex", index)
					return
				case task := <-m.messageChan:
					if err := m.sendMessage(task.namespaceUID, task.messageUID); err != nil {
						klog.Warnw("msg", "send message failed", "error", err, "namespaceUID", task.namespaceUID, "messageUID", task.messageUID)
						continue
					}
					klog.Debugw("msg", "send message success", "namespaceUID", task.namespaceUID, "messageUID", task.messageUID)
				}
			}
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
