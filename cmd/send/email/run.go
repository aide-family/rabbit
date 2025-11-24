package email

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kubeRegistry "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	kConfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
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
	c := kConfig.New(kConfig.WithSource(
		env.NewSource(),
		file.NewSource(flags.RabbitConfigPath),
	))
	if err := c.Load(); err != nil {
		flags.Helper.Errorw("msg", "load config failed", "error", err)
		return
	}
	if err := c.Scan(&bc); err != nil {
		flags.Helper.Errorw("msg", "scan config failed", "error", err)
		return
	}

	flags.applyToBootstrap(&bc)

	req, err := flags.parseRequestParams()
	if err != nil {
		flags.Helper.Errorw("msg", "parse request params failed", "error", err)
		return
	}
	var discovery connect.Registry
	switch registryType := bc.GetRegistryType(); registryType {
	case config.RegistryType_ETCD:
		etcdConfig := bc.GetEtcd()
		if pointer.IsNil(etcdConfig) {
			flags.Helper.Errorw("msg", "etcd config is not found")
			return
		}
		client, err := clientV3.New(clientV3.Config{
			Endpoints:   strutil.SplitSkipEmpty(etcdConfig.GetEndpoints(), ","),
			Username:    etcdConfig.GetUsername(),
			Password:    etcdConfig.GetPassword(),
			DialTimeout: 10 * time.Second,
		})
		if err != nil {
			flags.Helper.Errorw("msg", "etcd client initialization failed", "error", err)
			return
		}
		discovery = etcd.New(client)
	case config.RegistryType_KUBERNETES:
		kubeConfig := bc.GetKubernetes()
		if pointer.IsNil(kubeConfig) {
			flags.Helper.Errorw("msg", "kubernetes config is not found")
			return
		}
		kubeClient, err := connect.NewKubernetesClientSet(kubeConfig.GetKubeConfig())
		if err != nil {
			flags.Helper.Errorw("msg", "kubernetes client initialization failed", "error", err)
			return
		}
		discovery = kubeRegistry.NewRegistry(kubeClient, kubeConfig.GetNamespace())
	}

	clusterConfig := bc.GetCluster()
	clusterEndpoints := strutil.SplitSkipEmpty(clusterConfig.GetEndpoints(), ",")
	clusterProtocol := clusterConfig.GetProtocol()
	clusterTimeout := clusterConfig.GetTimeout().AsDuration()
	clusterName := clusterConfig.GetName()

	for _, clusterEndpoint := range clusterEndpoints {
		initConfig := connect.NewDefaultConfig(clusterName, clusterEndpoint, clusterTimeout)
		sender, err := NewSender(initConfig, clusterProtocol, bc.GetJwtToken(), discovery, flags.Helper)
		if err != nil {
			continue
		}

		reply, err := sender.SendEmail(context.Background(), req)
		if err != nil {
			flags.Helper.Warnw("msg", "send email failed", "cluster", clusterName, "protocol", clusterProtocol, "error", err)
			return
		}

		flags.Helper.Infow("msg", "send email success", "cluster", clusterName, "protocol", clusterProtocol, "reply", reply)
		return
	}
	// 没有可用的节点，退出
	flags.Helper.Error("no available nodes")
}

type Sender interface {
	SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error)
}

type sender struct {
	jwtToken string
	helper   *klog.Helper
	name     string
	protocol config.ClusterConfig_Protocol
	timeout  time.Duration
	call     func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error)
	close    func() error
}

// SendEmail implements Sender.
func (s *sender) SendEmail(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
	defer s.close()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ctx = metadata.NewClientContext(ctx, metadata.Metadata{
		cnst.MetadataGlobalKeyAuthorization: {s.jwtToken},
		cnst.MetadataGlobalKeyNamespace:     {flags.Namespace},
	})
	return s.call(ctx, in)
}

func NewSender(cluster connect.InitConfig, protocol config.ClusterConfig_Protocol, jwtToken string, discovery connect.Registry, helper *klog.Helper) (Sender, error) {
	name := cluster.GetName()
	newSender := &sender{
		jwtToken: jwtToken,
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
		connect.WithDiscovery(discovery),
	}
	switch newSender.protocol {
	case config.ClusterConfig_GRPC:
		grpcClient, err := connect.InitGRPCClient(cluster, opts...)
		if err != nil {
			helper.Errorw("msg", "cluster GRPC client initialization failed", "cluster", name, "protocol", protocol, "error", err)
			return nil, merr.ErrorInternalServer("failed to initialize GRPC client").WithCause(err)
		}
		newSender.close = grpcClient.Close
		newSender.call = func(ctx context.Context, in *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
			return apiv1.NewSenderClient(grpcClient).SendEmail(ctx, in)
		}
	case config.ClusterConfig_HTTP:
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
