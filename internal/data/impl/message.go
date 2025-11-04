package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aide-family/magicbox/safety"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

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
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewMessageBus(bc *conf.Bootstrap, d *data.Data, messageLogRepo repository.MessageLog, helper *klog.Helper) repository.MessageBus {
	eventBusConf := bc.GetEventBus()
	clustersConfig := bc.GetClusters()
	bus := &messageBusImpl{
		d:              d,
		messageLogRepo: messageLogRepo,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "component", "message_bus")),
		messageChan:    make(chan *messageTask, 1000), // 缓冲区大小
		senders:        safety.NewSyncMap(make(map[vobj.MessageType]repository.MessageSender)),
		stopChan:       make(chan struct{}),
		wg:             sync.WaitGroup{},
		workerCount:    int(eventBusConf.GetWorkerCount()),
		timeout:        eventBusConf.GetTimeout().AsDuration(),
		clusters:       make([]sender.Sender, 0, len(clustersConfig)),
	}

	// 注册发送器
	bus.senders.Set(vobj.MessageTypeEmail, sender.NewEmailSender(helper))
	bus.senders.Set(vobj.MessageTypeWebhook, sender.NewWebhookSender(helper))

	secret := bc.GetJwt().GetSecret()
	for _, clusterConfig := range clustersConfig {
		protocol := clusterConfig.GetProtocol()
		opts := []connect.InitOption{
			connect.WithProtocol(protocol.String()),
			connect.WithDiscovery(d.Registry()),
			connect.WithSecret(secret),
		}
		switch protocol {
		case config.ClusterConfig_GRPC:
			grpcClient, err := connect.InitGRPCClient(clusterConfig, opts...)
			if err != nil {
				helper.Errorw("msg", "create GRPC client failed", "error", err)
				continue
			}
			d.AppendClose("grpcClient."+clusterConfig.GetName(), func() error { return grpcClient.Close() })
			bus.clusters = append(bus.clusters, sender.NewClusterSender(grpcClient, nil, protocol))
		case config.ClusterConfig_HTTP:
			httpClient, err := connect.InitHTTPClient(clusterConfig, opts...)
			if err != nil {
				helper.Errorw("msg", "create HTTP client failed", "error", err)
				continue
			}
			d.AppendClose("httpClient."+clusterConfig.GetName(), func() error { return httpClient.Close() })
			bus.clusters = append(bus.clusters, sender.NewClusterSender(nil, httpClient, protocol))
		}
	}

	return bus
}

type messageTask struct {
	ctx     context.Context
	message *do.MessageLog
}

type messageBusImpl struct {
	d              *data.Data
	messageLogRepo repository.MessageLog
	helper         *klog.Helper
	messageChan    chan *messageTask
	senders        *safety.SyncMap[vobj.MessageType, repository.MessageSender]
	stopChan       chan struct{}
	wg             sync.WaitGroup
	workerCount    int // 工作协程数量,默认1个
	timeout        time.Duration

	clusters []sender.Sender
}

// start 启动后台处理goroutine
func (m *messageBusImpl) Start() {
	if m.workerCount <= 0 {
		m.workerCount = 1
	}
	m.wg.Add(m.workerCount)
	// 启动多个worker goroutine
	for i := 0; i < m.workerCount; i++ {
		go m.worker(i)
	}
}

// worker 处理消息的工作协程
func (m *messageBusImpl) worker(id int) {
	defer m.wg.Done()

	for {
		select {
		case task, ok := <-m.messageChan:
			if !ok {
				m.helper.Infow("msg", "message bus worker stopped", "worker", id)
				return
			}
			m.waitProcessMessage(task.ctx, task.message)
		case <-m.stopChan:
			m.helper.Infow("msg", "message bus worker stopped", "worker", id)
			return
		}
	}
}

func (m *messageBusImpl) waitProcessMessage(ctx context.Context, message *do.MessageLog) {
	req := &apiv1.SendMessageRequest{
		Uid: message.UID.Int64(),
	}
	for _, cluster := range m.clusters {
		reply, err := cluster.SendMessage(ctx, req)
		if err != nil {
			m.helper.Errorw("msg", "send message failed", "error", err, "uid", message.UID, "reply", reply)
			continue
		}
		return
	}
	if err := m.SendMessage(ctx, message); err != nil {
		m.helper.Errorw("msg", "send message failed", "error", err, "uid", message.UID)
	}
}

func (m *messageBusImpl) SendMessage(ctx context.Context, message *do.MessageLog) error {
	// 使用分布式锁控制并发，同一时间只能有一个服务或者协程能够访问此UID的消息
	namespace := middler.GetNamespace(ctx)
	// 在事务中使用 SELECT FOR UPDATE 获取分布式锁
	var newMessage *do.MessageLog
	err := m.d.BizDB(ctx, namespace).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = data.WithBizTransaction(ctx, tx, namespace)
		// 使用 SELECT FOR UPDATE 获取行锁，确保同一时间只有一个节点能处理该消息
		lockedMessage, err := m.messageLogRepo.GetMessageLogWithLock(ctx, message.UID)
		if err != nil {
			return err
		}

		// 如果消息已经发送或正在发送，直接返回
		if lockedMessage.Status.IsSent() || lockedMessage.Status.IsSending() {
			return nil
		}

		// 使用 CAS 操作原子性地更新状态为发送中
		// 只有当前状态为待处理或失败时才更新为发送中
		result, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, message.UID, vobj.MessageStatusPending, vobj.MessageStatusSending)
		if err != nil {
			return err
		}
		if !result {
			return nil
		}

		newMessage = lockedMessage
		// 更新状态为发送中
		newMessage.Status = vobj.MessageStatusSending
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
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	senderType := message.Type
	sender, ok := m.senders.Get(senderType)
	if !ok {
		m.helper.Errorw("msg", "sender not found", "type", senderType, "uid", message.UID)
		if err := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusFailed); err != nil {
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
func (m *messageBusImpl) AppendMessage(ctx context.Context, message *do.MessageLog) error {
	// 将消息放入channel异步处理
	select {
	case m.messageChan <- &messageTask{ctx: safety.CopyValueCtx(ctx), message: message}:
		return nil
	default:
		// channel满了,返回错误
		m.helper.Errorw("msg", "message channel is full", "uid", message.UID)
		return fmt.Errorf("message channel is full")
	}
}

// Stop 停止事件总线
func (m *messageBusImpl) Stop() {
	close(m.stopChan)
	m.wg.Wait()
	close(m.messageChan)
}
