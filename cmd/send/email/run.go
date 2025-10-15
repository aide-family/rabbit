package email

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/internal/conf"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/merr"
)

func run(cmd *cobra.Command, args []string) {
	var bc conf.Bootstrap
	if err := load.Load(flags.configPath, &bc); err != nil {
		panic(err)
	}

	flags.applyToBootstrap(&bc)
	logger, err := log.NewLogger(stdio.LoggerDriver())
	if err != nil {
		panic(err)
	}
	helper := klog.NewHelper(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &apiv1.SendEmailRequest{
		Namespace:   flags.Namespace,
		Subject:     flags.Subject,
		Body:        flags.Body,
		To:          flags.To,
		Cc:          flags.Cc,
		ContentType: flags.ContentType,
		Headers:     flags.Headers,
	}
	for _, cluster := range bc.GetClusters() {
		sender, err := NewSender(cluster)
		if err != nil {
			helper.Warn(err)
			continue
		}
		reply, err := sender.SendEmail(ctx, req)
		if err != nil {
			helper.Warn(err)
			continue
		}
		helper.Info(reply)
		return
	}
	// 没有可用的节点，退出
	helper.Error("no available nodes")
}

type Sender interface {
	SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error)
}

type sender struct {
	call func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error)
}

// SendEmail implements Sender.
func (s *sender) SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
	return s.call(ctx, in)
}

func NewSender(cluster *conf.Cluster) (Sender, error) {
	switch cluster.GetProtocol() {
	case conf.Cluster_GRPC:
		grpcClient, err := connect.InitGRPCClient(cluster)
		if err != nil {
			return nil, err
		}
		emailService := apiv1.NewEmailClient(grpcClient)
		return &sender{
			call: func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
				return emailService.SendEmail(ctx, in)
			},
		}, nil
	case conf.Cluster_HTTP:
		httpClient, err := connect.InitHTTPClient(cluster)
		if err != nil {
			return nil, err
		}
		emailService := apiv1.NewEmailHTTPClient(httpClient)
		return &sender{
			call: func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
				return emailService.SendEmail(ctx, in)
			},
		}, nil
	default:
		return nil, merr.ErrorInternalServer("unknown protocol")
	}
}
