package connect

import (
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	ProtocolHTTP = "HTTP"
	ProtocolGRPC = "GRPC"
)

type InitConfig interface {
	GetName() string
	GetEndpoint() string
	GetTimeout() *durationpb.Duration
}

type initConfig struct {
	name        string
	endpoint    string
	protocol    string
	timeout     time.Duration
	nodeVersion string
	secret      string
	discovery   registry.Discovery
}

func NewInitConfig(config InitConfig, opts ...InitOption) *initConfig {
	cfg := &initConfig{
		name:     config.GetName(),
		endpoint: config.GetEndpoint(),
		protocol: ProtocolGRPC,
		timeout:  config.GetTimeout().AsDuration(),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type InitOption func(*initConfig)

func WithNodeVersion(version string) InitOption {
	return func(cfg *initConfig) {
		cfg.nodeVersion = version
	}
}

func WithDiscovery(discovery registry.Discovery) InitOption {
	return func(cfg *initConfig) {
		cfg.discovery = discovery
	}
}

func WithSecret(secret string) InitOption {
	return func(cfg *initConfig) {
		cfg.secret = secret
	}
}

func WithProtocol(protocol string) InitOption {
	return func(cfg *initConfig) {
		cfg.protocol = protocol
	}
}
