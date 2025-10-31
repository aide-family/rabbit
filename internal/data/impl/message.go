package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aide-family/magicbox/safety"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/sender"
	"github.com/aide-family/rabbit/pkg/config"
)

func NewMessageBus(bc *conf.Bootstrap, d *data.Data, messageLogRepo repository.MessageLog, helper *klog.Helper) repository.MessageBus {
	eventBusConf := bc.GetEventBus()
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
		clusters:       bc.GetClusters(),
	}

	// 注册发送器
	bus.senders.Set(vobj.MessageTypeEmail, sender.NewEmailSender(helper))
	bus.senders.Set(vobj.MessageTypeWebhook, sender.NewWebhookSender(helper))

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

	clusters []*config.Cluster
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
}

func (m *messageBusImpl) SendMessage(ctx context.Context, message *do.MessageLog) error {
	m.processMessage(ctx, message)
	return nil
}

// processMessage 处理消息
func (m *messageBusImpl) processMessage(ctx context.Context, message *do.MessageLog) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	// 更新状态为发送中
	if err := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusSending); err != nil {
		m.helper.Errorw("msg", "update message status to sending failed", "error", err, "uid", message.UID)
		return
	}

	senderType := message.Type
	sender, ok := m.senders.Get(senderType)
	if !ok {
		m.helper.Errorw("msg", "sender not found", "type", senderType, "uid", message.UID)
		if err := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusFailed); err != nil {
			m.helper.Errorw("msg", "update message status to failed failed", "error", err, "uid", message.UID)
		}
		return
	}

	// 发送消息
	if err := sender.Send(ctx, message); err != nil {
		m.helper.Errorw("msg", "send message failed", "error", err, "uid", message.UID, "type", senderType)
		if updateErr := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusFailed); updateErr != nil {
			m.helper.Errorw("msg", "update message status to failed failed", "error", updateErr, "uid", message.UID)
		}
		return
	}

	// 更新状态为已发送
	if err := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusSent); err != nil {
		m.helper.Errorw("msg", "update message status to sent failed", "error", err, "uid", message.UID)
		return
	}

	m.helper.Infow("msg", "message sent successfully", "uid", message.UID, "type", senderType)
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
