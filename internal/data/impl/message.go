package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/sender"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewMessageBus(
	bc *conf.Bootstrap,
	d *data.Data,
	transactionRepo repository.Transaction,
	messageLogRepo repository.MessageLog,
	helper *klog.Helper,
) repository.MessageBus {
	eventBusConf := bc.GetEventBus()
	clusterConfig := bc.GetCluster()
	clusterEndpoints := strutil.SplitSkipEmpty(clusterConfig.GetEndpoints(), ",")
	clusterProtocol := clusterConfig.GetProtocol()
	clusterTimeout := clusterConfig.GetTimeout().AsDuration()
	clusterName := clusterConfig.GetName()
	bus := &messageBusImpl{
		d:               d,
		transactionRepo: transactionRepo,
		messageLogRepo:  messageLogRepo,
		helper:          klog.NewHelper(klog.With(helper.Logger(), "component", "message_bus")),
		messageChan:     make(chan *messageTask, eventBusConf.GetBufferSize()),
		senders:         safety.NewSyncMap(make(map[vobj.MessageType]repository.MessageSender)),
		stopChan:        make(chan struct{}),
		wg:              sync.WaitGroup{},
		workerCount:     int(eventBusConf.GetWorkerCount()),
		timeout:         eventBusConf.GetTimeout().AsDuration(),
		clusters:        make([]sender.Sender, 0, len(clusterEndpoints)),
	}

	// 注册发送器
	bus.senders.Set(vobj.MessageTypeEmail, sender.NewEmailSender(helper))
	bus.senders.Set(vobj.MessageTypeWebhook, sender.NewWebhookSender(helper))

	safety.Go(context.Background(), "message_bus_init_clusters", func(ctx context.Context) error {
		time.Sleep(3 * time.Second)
		bus.initClusters(clusterEndpoints, clusterProtocol, clusterName, clusterTimeout)
		return nil
	}, helper.Logger())

	return bus
}

type messageTask struct {
	ctx        context.Context
	messageUID snowflake.ID
}

type messageBusImpl struct {
	d               *data.Data
	transactionRepo repository.Transaction
	messageLogRepo  repository.MessageLog
	helper          *klog.Helper
	messageChan     chan *messageTask
	senders         *safety.SyncMap[vobj.MessageType, repository.MessageSender]
	stopChan        chan struct{}
	wg              sync.WaitGroup
	workerCount     int // 工作协程数量,默认1个
	timeout         time.Duration

	clusters []sender.Sender
}

// start 启动后台处理goroutine
func (m *messageBusImpl) Start(ctx context.Context) {
	if m.workerCount <= 0 {
		m.workerCount = 1
	}

	// 启动多个worker goroutine
	for i := 0; i < m.workerCount; i++ {
		m.wg.Add(1)
		safety.Go(ctx, fmt.Sprintf("message_bus_worker_%d", i), func(ctx context.Context) error {
			defer m.wg.Done()
			m.worker(ctx, i)
			return nil
		}, m.helper.Logger())
	}
}

// Stop 停止事件总线
func (m *messageBusImpl) Stop(ctx context.Context) {
	select {
	case <-ctx.Done():
		m.helper.Warnw("msg", "message bus stopped by context done")
		return
	case <-m.stopChan:
		m.helper.Warnw("msg", "message bus stopped by stop channel")
		return
	default:
		close(m.stopChan)
		m.wg.Wait()
		close(m.messageChan)
		m.helper.Infow("msg", "message bus stopped")
	}
}

// worker 处理消息的工作协程
func (m *messageBusImpl) worker(ctx context.Context, id int) {
	for {
		select {
		case task, ok := <-m.messageChan:
			if !ok {
				m.helper.Infow("msg", "message bus worker stopped by message channel closed", "worker", id)
				return
			}
			m.waitProcessMessage(task.ctx, task.messageUID)
		case <-m.stopChan:
			m.helper.Infow("msg", "message bus worker stopped by stop channel", "worker", id)
			return
		case <-ctx.Done():
			m.helper.Infow("msg", "message bus worker stopped by context done", "worker", id)
			return
		}
	}
}

func (m *messageBusImpl) waitProcessMessage(ctx context.Context, messageUID snowflake.ID) {
	req := &apiv1.SendMessageRequest{
		Uid: messageUID.Int64(),
	}
	for _, cluster := range m.clusters {
		reply, err := cluster.SendMessage(ctx, req)
		if err != nil {
			m.helper.Errorw("msg", "send message failed", "error", err, "uid", messageUID, "reply", reply)
			continue
		}
		return
	}
	if err := m.SendMessage(ctx, messageUID); err != nil {
		m.helper.Errorw("msg", "send message failed", "error", err, "uid", messageUID)
	}
}

func (m *messageBusImpl) SendMessage(ctx context.Context, messageUID snowflake.ID) error {
	// 在事务中使用 SELECT FOR UPDATE 获取分布式锁
	var newMessage *do.MessageLog
	err := m.transactionRepo.Transaction(ctx, func(transactionCtx context.Context) error {
		// 使用 SELECT FOR UPDATE 获取行锁，确保同一时间只有一个节点能处理该消息
		lockedMessage, err := m.messageLogRepo.GetMessageLogWithLock(transactionCtx, messageUID)
		if err != nil {
			if merr.IsNotFound(err) {
				return nil
			}
			return merr.ErrorInternal("get message log with lock failed").WithCause(err)
		}

		// 如果消息已经发送或正在发送，直接返回
		if lockedMessage.Status.IsSent() || lockedMessage.Status.IsSending() {
			return nil
		}

		// 使用 CAS 操作原子性地更新状态为发送中
		// 只有当前状态为待处理或失败时才更新为发送中
		result, err := m.messageLogRepo.UpdateMessageLogStatusIf(transactionCtx, messageUID, vobj.MessageStatusPending, vobj.MessageStatusSending)
		if err != nil {
			return merr.ErrorInternal("update message status to sending failed").WithCause(err)
		}
		if !result {
			return nil
		}

		newMessage = lockedMessage
		return nil
	})
	if err != nil {
		return err
	}

	// 如果消息已经被处理或者状态更新失败，直接返回
	if newMessage == nil {
		return nil
	}

	// 在事务外处理消息发送（发送操作可能需要较长时间，不应该在数据库事务中执行）
	return m.processMessage(ctx, newMessage)
}

// processMessage 处理消息
func (m *messageBusImpl) processMessage(ctx context.Context, message *do.MessageLog) error {
	if message.Status.IsSent() || message.Status.IsSending() {
		return nil
	}
	message.Status = vobj.MessageStatusSending
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	senderType := message.Type
	sender, ok := m.senders.Get(senderType)
	if !ok {
		m.helper.Errorw("msg", "sender not found", "type", senderType, "uid", message.UID)
		if _, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, message.UID, vobj.MessageStatusSending, vobj.MessageStatusFailed); err != nil {
			m.helper.Errorw("msg", "update message status to failed failed", "error", err, "uid", message.UID)
		}
		return merr.ErrorParams("sender not found")
	}

	// 发送消息
	if err := sender.Send(ctx, message); err != nil {
		m.helper.Errorw("msg", "send message failed", "error", err, "uid", message.UID, "type", senderType)
		success, updateErr := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, message.UID, vobj.MessageStatusSending, vobj.MessageStatusFailed)
		if updateErr != nil {
			m.helper.Errorw("msg", "update message status to failed failed", "error", updateErr, "uid", message.UID)
		}
		if !success {
			m.helper.Warnw("msg", "message status is not sending, message sent failed", "uid", message.UID, "type", senderType)
		}
		return merr.ErrorInternal("send message failed").WithCause(err)
	}

	// 更新状态为已发送
	success, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, message.UID, vobj.MessageStatusSending, vobj.MessageStatusSent)
	if err != nil {
		m.helper.Errorw("msg", "update message status to sent failed", "error", err, "uid", message.UID)
		return merr.ErrorInternal("update message status to sent failed")
	}
	if !success {
		m.helper.Warnw("msg", "message status is not sending, message sent successfully", "uid", message.UID, "type", senderType)
	}
	return nil
}

// AppendMessage implements repository.MessageBus.
func (m *messageBusImpl) AppendMessage(ctx context.Context, messageUID snowflake.ID) error {
	// 将消息放入channel异步处理
	select {
	case m.messageChan <- &messageTask{ctx: safety.CopyValueCtx(ctx), messageUID: messageUID}:
		return nil
	default:
		// channel满了,返回错误
		m.helper.Errorw("msg", "message channel is full", "uid", messageUID)
		return merr.ErrorInternal("message channel is full")
	}
}

func (m *messageBusImpl) initClusters(clusterEndpoints []string, clusterProtocol config.ClusterConfig_Protocol, clusterName string, clusterTimeout time.Duration) {
	for _, clusterEndpoint := range clusterEndpoints {
		opts := []connect.InitOption{
			connect.WithProtocol(clusterProtocol.String()),
			connect.WithDiscovery(m.d.Registry()),
		}
		initConfig := connect.NewDefaultConfig(clusterName, clusterEndpoint, clusterTimeout)
		switch clusterProtocol {
		case config.ClusterConfig_GRPC:
			grpcClient, err := connect.InitGRPCClient(initConfig, opts...)
			if err != nil {
				m.helper.Errorw("msg", "create GRPC client failed", "endpoint", clusterEndpoint, "error", err)
				continue
			}
			m.d.AppendClose("grpcClient", func() error { return grpcClient.Close() })
			m.clusters = append(m.clusters, sender.NewClusterSender(grpcClient, nil, clusterProtocol))
		case config.ClusterConfig_HTTP:
			httpClient, err := connect.InitHTTPClient(initConfig, opts...)
			if err != nil {
				m.helper.Errorw("msg", "create HTTP client failed", "error", err)
				continue
			}
			m.d.AppendClose("httpClient", func() error { return httpClient.Close() })
			m.clusters = append(m.clusters, sender.NewClusterSender(nil, httpClient, clusterProtocol))
		}
	}
}
