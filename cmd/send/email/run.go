package email

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/spf13/cobra"
	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/aide-family/rabbit/cmd"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/merr"
)

func run(_ *cobra.Command, _ []string) {
	flags.GlobalFlags = cmd.GetGlobalFlags()
	var bc config.ClientConfig
	if err := load.Load(load.ExpandHomeDir(flags.RabbitConfigPath), &bc); err != nil {
		panic(err)
	}

	flags.applyToBootstrap(&bc)

	req, err := flags.parseRequestParams()
	if err != nil {
		flags.Helper.Errorw("msg", "parse request params failed", "error", err)
		return
	}
	var discovery registry.Discovery
	if pointer.IsNotNil(bc.GetEtcd()) {
		client, err := clientV3.New(clientV3.Config{
			Endpoints:   bc.GetEtcd().GetEndpoints(),
			Username:    bc.GetEtcd().GetUsername(),
			Password:    bc.GetEtcd().GetPassword(),
			DialTimeout: 10 * time.Second,
		})
		if err != nil {
			flags.Helper.Errorw("msg", "etcd client initialization failed", "error", err)
			return
		}
		discovery = etcd.New(client)
	}

	for _, cluster := range bc.GetClusters() {
		sender, err := NewSender(cluster, discovery, bc.GetJwtToken(), flags.Helper)
		if err != nil {
			continue
		}
		name, protocol := cluster.GetName(), cluster.GetProtocol()
		reply, err := sender.SendEmail(context.Background(), req)
		if err != nil {
			flags.Helper.Warnw("msg", "send email failed", "cluster", name, "protocol", protocol, "error", err)
			break
		}

		flags.Helper.Infow("msg", "send email success", "cluster", name, "protocol", protocol, "reply", reply)
		return
	}
	// 没有可用的节点，退出
	flags.Helper.Error("no available nodes")
}

type Sender interface {
	SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error)
}

type sender struct {
	helper   *klog.Helper
	name     string
	protocol config.Cluster_Protocol
	timeout  time.Duration
	call     func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error)
	close    func() error
}

// SendEmail implements Sender.
func (s *sender) SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
	defer s.close()
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	ctx = metadata.AppendToClientContext(ctx, cnst.MetadataGlobalKeyNamespace, flags.Namespace)
	return s.call(ctx, in)
}

func NewSender(cluster *config.Cluster, discovery registry.Discovery, token string, helper *klog.Helper) (Sender, error) {
	name, protocol, secret := cluster.GetName(), cluster.GetProtocol(), cluster.GetSecret()
	newSender := &sender{
		helper:   helper,
		name:     name,
		protocol: protocol,
		timeout:  cluster.GetTimeout().AsDuration(),
		close:    func() error { return nil },
		call: func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
			helper.Errorw("msg", "unknown protocol", "cluster", name, "protocol", protocol)
			return nil, merr.ErrorInternalServer("cluster %s unknown protocol %s", name, protocol)
		},
	}
	opts := []connect.InitOption{
		connect.WithProtocol(newSender.protocol.String()),
		connect.WithSecret(secret),
		connect.WithToken(token),
		connect.WithDiscovery(discovery),
	}
	switch newSender.protocol {
	case config.Cluster_GRPC:
		grpcClient, err := connect.InitGRPCClient(cluster, opts...)
		if err != nil {
			helper.Errorw("msg", "cluster GRPC client initialization failed", "cluster", name, "protocol", protocol, "error", err)
			return nil, merr.ErrorInternalServer("failed to initialize GRPC client").WithCause(err)
		}
		newSender.close = grpcClient.Close
		newSender.call = func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
			return apiv1.NewSenderClient(grpcClient).SendEmail(ctx, in)
		}
	case config.Cluster_HTTP:
		httpClient, err := connect.InitHTTPClient(cluster, opts...)
		if err != nil {
			helper.Errorw("msg", "cluster HTTP client initialization failed", "cluster", name, "protocol", protocol, "error", err)
			return nil, merr.ErrorInternalServer("failed to initialize HTTP client").WithCause(err)
		}
		newSender.close = httpClient.Close
		newSender.call = func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
			return apiv1.NewSenderHTTPClient(httpClient).SendEmail(ctx, in)
		}
	}
	return newSender, nil
}
