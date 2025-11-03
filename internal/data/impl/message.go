package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aide-family/magicbox/safety"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"

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

func newClusterSender(conn *grpc.ClientConn, http *http.Client, protocol config.ClusterConfig_Protocol) sender.Sender {
	return &clusterSender{
		conn:     conn,
		http:     http,
		protocol: protocol,
	}
}

type clusterSender struct {
	conn     *grpc.ClientConn
	http     *http.Client
	protocol config.ClusterConfig_Protocol
}

func (c *clusterSender) SendMessage(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error) {
	switch c.protocol {
	case config.ClusterConfig_GRPC:
		return apiv1.NewSenderClient(c.conn).SendMessage(ctx, req)
	case config.ClusterConfig_HTTP:
		return apiv1.NewSenderHTTPClient(c.http).SendMessage(ctx, req)
	}
	return nil, nil
}

func (c *clusterSender) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	if c.http != nil {
		return c.http.Close()
	}
	return nil
}

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
			bus.clusters = append(bus.clusters, newClusterSender(grpcClient, nil, protocol))
		case config.ClusterConfig_HTTP:
			httpClient, err := connect.InitHTTPClient(clusterConfig, opts...)
			if err != nil {
				helper.Errorw("msg", "create HTTP client failed", "error", err)
				continue
			}
			d.AppendClose("httpClient."+clusterConfig.GetName(), func() error { return httpClient.Close() })
			bus.clusters = append(bus.clusters, newClusterSender(nil, httpClient, protocol))
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
	newMessage, err := m.messageLogRepo.GetMessageLog(ctx, message.UID)
	if err != nil {
		return err
	}
	if newMessage.Status.IsSent() || newMessage.Status.IsSending() {
		return nil
	}
	return m.processMessage(ctx, newMessage)
}

// processMessage 处理消息
func (m *messageBusImpl) processMessage(ctx context.Context, message *do.MessageLog) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	// 更新状态为发送中
	if err := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusSending); err != nil {
		m.helper.Errorw("msg", "update message status to sending failed", "error", err, "uid", message.UID)
		return merr.ErrorInternal("update message status to sending failed")
	}

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
		if updateErr := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusFailed); updateErr != nil {
			m.helper.Errorw("msg", "update message status to failed failed", "error", updateErr, "uid", message.UID)
		}
		return merr.ErrorInternal("send message failed")
	}

	// 更新状态为已发送
	if err := m.messageLogRepo.UpdateMessageLogStatus(ctx, message.UID, vobj.MessageStatusSent); err != nil {
		m.helper.Errorw("msg", "update message status to sent failed", "error", err, "uid", message.UID)
		return merr.ErrorInternal("update message status to sent failed")
	}

	m.helper.Infow("msg", "message sent successfully", "uid", message.UID, "type", senderType)
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
