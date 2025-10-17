package email

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/merr"
)

func run(cmd *cobra.Command, args []string) {
	var bc config.ClientConfig
	if err := load.Load(load.ExpandHomeDir(flags.configPath), &bc); err != nil {
		panic(err)
	}

	flags.applyToBootstrap(&bc)
	logger, err := log.NewLogger(stdio.LoggerDriver())
	if err != nil {
		panic(err)
	}
	helper := klog.NewHelper(logger)

	req, err := flags.parseRequestParams()
	if err != nil {
		helper.Errorw("msg", "parse request params failed", "error", err)
		return
	}
	ctx := context.Background()
	for _, cluster := range bc.GetClusters() {
		sender, err := NewSender(cluster, bc.GetJwtToken(), helper)
		if err != nil {
			continue
		}
		name, protocol := cluster.GetName(), cluster.GetProtocol()
		reply, err := sender.SendEmail(ctx, req)
		if err != nil {
			helper.Warnw("msg", "send email failed", "cluster", name, "protocol", protocol, "error", err)
			continue
		}

		helper.Infow("msg", "send email success", "cluster", name, "protocol", protocol, "reply", reply)
		return
	}
	// 没有可用的节点，退出
	helper.Error("no available nodes")
}

type Sender interface {
	SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error)
}

type sender struct {
	helper   *klog.Helper
	name     string
	protocol config.Cluster_Protocol
	timeout  time.Duration
	call     func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error)
	close    func() error
}

// SendEmail implements Sender.
func (s *sender) SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
	defer s.close()
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.call(ctx, in)
}

func NewSender(cluster *config.Cluster, token string, helper *klog.Helper) (Sender, error) {
	name, protocol, secret := cluster.GetName(), cluster.GetProtocol(), cluster.GetSecret()
	newSender := &sender{
		helper:   helper,
		name:     name,
		protocol: protocol,
		timeout:  cluster.GetTimeout().AsDuration(),
		close:    func() error { return nil },
		call: func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
			helper.Errorw("msg", "unknown protocol", "cluster", name, "protocol", protocol)
			return nil, merr.ErrorInternalServer("cluster %s unknown protocol %s", name, protocol)
		},
	}
	opts := []connect.InitOption{
		connect.WithProtocol(newSender.protocol.String()),
		connect.WithSecret(secret),
		connect.WithToken(token),
	}
	switch newSender.protocol {
	case config.Cluster_GRPC:
		grpcClient, err := connect.InitGRPCClient(cluster, opts...)
		if err != nil {
			helper.Errorw("msg", "cluster GRPC client initialization failed", "cluster", name, "protocol", protocol, "error", err)
			return nil, merr.ErrorInternalServer("failed to initialize GRPC client").WithCause(err)
		}
		newSender.close = grpcClient.Close
		newSender.call = func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
			return apiv1.NewEmailClient(grpcClient).SendEmail(ctx, in)
		}
	case config.Cluster_HTTP:
		httpClient, err := connect.InitHTTPClient(cluster, opts...)
		if err != nil {
			helper.Errorw("msg", "cluster HTTP client initialization failed", "cluster", name, "protocol", protocol, "error", err)
			return nil, merr.ErrorInternalServer("failed to initialize HTTP client").WithCause(err)
		}
		newSender.close = httpClient.Close
		newSender.call = func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
			return apiv1.NewEmailHTTPClient(httpClient).SendEmail(ctx, in)
		}
	}
	return newSender, nil
}
