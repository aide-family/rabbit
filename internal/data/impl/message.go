package impl

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
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

func NewMessageRepository(
	bc *conf.Bootstrap,
	d *data.Data,
	transactionRepo repository.Transaction,
	messageLogRepo repository.MessageLog,
	helper *klog.Helper,
) repository.Message {
	jobCoreConf := bc.GetJobCore()
	clusterConfig := bc.GetCluster()
	clusterEndpoints := strutil.SplitSkipEmpty(clusterConfig.GetEndpoints(), ",")
	messageRepo := &messageRepositoryImpl{
		d:               d,
		bc:              bc,
		transactionRepo: transactionRepo,
		messageLogRepo:  messageLogRepo,
		helper:          klog.NewHelper(klog.With(helper.Logger(), "impl", "message")),
		messageChan:     make(chan *messageTask, jobCoreConf.GetBufferSize()),
		senders:         safety.NewSyncMap(make(map[vobj.MessageType]repository.MessageSender)),
		stopChan:        make(chan struct{}),
		wg:              sync.WaitGroup{},
		workerTotal:     int(jobCoreConf.GetWorkerTotal()),
		timeout:         jobCoreConf.GetTimeout().AsDuration(),
		clusters:        make([]sender.Sender, 0, len(clusterEndpoints)),
	}

	// 注册发送器
	messageRepo.registerSenders(sender.NewEmailSender(helper), sender.NewWebhookSender(helper))

	messageRepo.Start(context.Background())

	messageRepo.d.AppendClose("messageRepository", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return messageRepo.Stop(ctx)
	})

	return messageRepo
}

type messageTask struct {
	ctx        context.Context
	messageUID snowflake.ID
}

type messageRepositoryImpl struct {
	d               *data.Data
	bc              *conf.Bootstrap
	transactionRepo repository.Transaction
	messageLogRepo  repository.MessageLog
	helper          *klog.Helper
	messageChan     chan *messageTask
	senders         *safety.SyncMap[vobj.MessageType, repository.MessageSender]
	stopChan        chan struct{}
	wg              sync.WaitGroup
	workerTotal     int // 工作协程数量,默认1个
	timeout         time.Duration

	clusters        []sender.Sender
	clusterInitOnce sync.Once
}

// start 启动后台处理goroutine
func (m *messageRepositoryImpl) Start(ctx context.Context) error {
	if m.workerTotal <= 0 {
		m.workerTotal = 1
	}

	// 启动多个worker goroutine
	for workerID := 0; workerID < m.workerTotal; workerID++ {
		m.worker(ctx, workerID)
	}
	return nil
}

// Stop 停止事件总线
func (m *messageRepositoryImpl) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		m.helper.Debug("msg", "message bus stopped by context done")
		return nil
	case <-m.stopChan:
		m.helper.Debug("msg", "message bus stopped by stop channel")
		return nil
	default:
		close(m.stopChan)
		m.wg.Wait()
		close(m.messageChan)
		m.helper.Debug("msg", "message bus stopped")
		return nil
	}
}

// worker 处理消息的工作协程
func (m *messageRepositoryImpl) worker(ctx context.Context, workerID int) {
	m.wg.Go(func() {
		for {
			select {
			case task, ok := <-m.messageChan:
				if !ok {
					m.helper.Debugw("msg", "message bus worker stopped by message channel closed", "worker", workerID)
					return
				}
				m.waitProcessMessage(task.ctx, task.messageUID)
			case <-m.stopChan:
				m.helper.Debugw("msg", "message bus worker stopped by stop channel", "worker", workerID)
				return
			case <-ctx.Done():
				m.helper.Debugw("msg", "message bus worker stopped by context done", "worker", workerID)
				return
			}
		}
	})
}

func (m *messageRepositoryImpl) waitProcessMessage(ctx context.Context, messageUID snowflake.ID) {
	req := &apiv1.JobSendMessageRequest{
		Uid: messageUID.Int64(),
	}
	// notice: 没有使用外部存储，不允许使用集群模式， 避免消息无法共享到其他节点
	if m.d.UseDatabase() {
		m.initClusters()
		rand.Shuffle(len(m.clusters), func(i, j int) {
			m.clusters[i], m.clusters[j] = m.clusters[j], m.clusters[i]
		})

		// 按打乱后的顺序尝试发送，失败则重试下一个节点
		for _, cluster := range m.clusters {
			reply, err := cluster.SendMessage(ctx, req)
			if err != nil {
				m.helper.Errorw("msg", "send message failed", "error", err, "uid", messageUID, "reply", reply, "cluster", cluster)
				continue
			}
			return
		}
		m.helper.Debugw("msg", "no cluster available to send message, use local node", "uid", messageUID)
	}
	// 如果未启用数据库或者全部节点都失败，则直接使用当前节点发送消息
	if err := m.SendMessage(ctx, messageUID); err != nil {
		m.helper.Errorw("msg", "send message failed", "error", err, "uid", messageUID)
	}
}

func (m *messageRepositoryImpl) SendMessage(ctx context.Context, messageUID snowflake.ID) error {
	// 在事务中使用 SELECT FOR UPDATE 获取分布式锁
	var newMessage *bo.MessageLogItemBo
	err := m.transactionRepo.Transaction(ctx, func(transactionCtx context.Context) error {
		// 使用 SELECT FOR UPDATE 获取行锁，确保同一时间只有一个节点能处理该消息
		lockedMessage, err := m.messageLogRepo.GetMessageLogWithLock(transactionCtx, messageUID)
		if err != nil {
			if merr.IsNotFound(err) {
				m.helper.Debugw("msg", "message log not found", "uid", messageUID)
				return nil
			}
			return merr.ErrorInternal("get message log with lock failed").WithCause(err)
		}

		// 如果消息已经发送或正在发送，直接返回
		if lockedMessage.Status.IsSent() || lockedMessage.Status.IsSending() {
			m.helper.Debugw("msg", "message status is not sent or sending, skip update status", "uid", messageUID, "status", lockedMessage.Status)
			return nil
		}

		// 使用 CAS 操作原子性地更新状态为发送中
		// 只有当前状态为待处理或失败时才更新为发送中
		result, err := m.messageLogRepo.UpdateMessageLogStatusIf(transactionCtx, messageUID, vobj.MessageStatusPending, vobj.MessageStatusSending)
		if err != nil {
			return merr.ErrorInternal("update message status to sending failed").WithCause(err)
		}
		if !result {
			m.helper.Debugw("msg", "message already processed or status update failed", "uid", messageUID)
			return nil
		}

		newMessage = bo.NewMessageLogItemBo(lockedMessage)
		return nil
	})
	if err != nil {
		m.helper.Errorw("msg", "transaction failed", "error", err, "uid", messageUID)
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
func (m *messageRepositoryImpl) processMessage(ctx context.Context, message *bo.MessageLogItemBo) error {
	if message.Status.IsSent() || message.Status.IsSending() {
		m.helper.Debugw("msg", "message already sent or sending", "uid", message.UID)
		return nil
	}
	message.Status = vobj.MessageStatusSending
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	senderType := message.Type
	sender, ok := m.senders.Get(senderType)
	if !ok {
		m.helper.Debugw("msg", "sender not found", "type", senderType, "uid", message.UID)
		if _, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, message.UID, vobj.MessageStatusSending, vobj.MessageStatusFailed); err != nil {
			m.helper.Errorw("msg", "update message status to failed failed", "error", err, "uid", message.UID)
		}
		return merr.ErrorParams("sender not supported")
	}

	// 发送消息
	if err := sender.Send(ctx, message); err != nil {
		m.helper.Errorw("msg", "send message failed", "error", err, "uid", message.UID, "type", senderType)
		success, updateErr := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, message.UID, vobj.MessageStatusSending, vobj.MessageStatusFailed)
		if updateErr != nil {
			m.helper.Errorw("msg", "update message status to failed failed", "error", updateErr, "uid", message.UID)
		}
		if !success {
			m.helper.Debugw("msg", "message status is not sending, message sent failed", "uid", message.UID, "type", senderType)
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
		m.helper.Debugw("msg", "message status is not sending, message sent successfully", "uid", message.UID, "type", senderType)
	}
	return nil
}

// AppendMessage implements repository.Message.
func (m *messageRepositoryImpl) AppendMessage(ctx context.Context, messageUID snowflake.ID) error {
	// 将消息放入channel异步处理
	select {
	case m.messageChan <- &messageTask{ctx: safety.CopyValueCtx(ctx), messageUID: messageUID}:
		m.helper.Debugw("msg", "message appended to channel", "uid", messageUID)
		return nil
	default:
		// channel满了,返回错误
		m.helper.Debugw("msg", "message channel is full", "uid", messageUID)
		return merr.ErrorInternal("message channel is full")
	}
}

func (m *messageRepositoryImpl) initClusters() {
	m.clusterInitOnce.Do(func() {
		clusterConfig := m.bc.GetCluster()
		clusterEndpoints := strutil.SplitSkipEmpty(clusterConfig.GetEndpoints(), ",")
		clusterTimeout := clusterConfig.GetTimeout().AsDuration()
		clusterName := clusterConfig.GetName()
		for _, clusterEndpoint := range clusterEndpoints {
			opts := []connect.InitOption{
				connect.WithProtocol(config.ClusterConfig_GRPC.String()),
				connect.WithDiscovery(m.d.Registry()),
			}
			initConfig := connect.NewDefaultConfig(clusterName, clusterEndpoint, clusterTimeout)
			grpcClient, err := connect.InitGRPCClient(initConfig, opts...)
			if err != nil {
				m.helper.Errorw("msg", "create GRPC client failed", "endpoint", clusterEndpoint, "error", err)
				continue
			}
			m.d.AppendClose("grpcClient", func() error { return grpcClient.Close() })
			m.clusters = append(m.clusters, sender.NewClusterSender(grpcClient))
		}
	})
}

func (m *messageRepositoryImpl) registerSenders(senders ...repository.MessageSender) {
	for _, sender := range senders {
		m.senders.Set(sender.Type(), sender)
	}
}
