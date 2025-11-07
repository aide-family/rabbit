package connect

import (
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"google.golang.org/protobuf/types/known/durationpb"

	rabbitMiddler "github.com/aide-family/rabbit/pkg/middler"
)

func init() {
	selector.SetGlobalSelector(p2c.NewBuilder())
}

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
	claim       *rabbitMiddler.JwtClaims
}

func NewInitConfig(config InitConfig, opts ...InitOption) (*initConfig, error) {
	cfg := &initConfig{
		name:     config.GetName(),
		endpoint: config.GetEndpoint(),
		protocol: ProtocolGRPC,
		timeout:  config.GetTimeout().AsDuration(),
		claim:    &rabbitMiddler.JwtClaims{},
	}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

type InitOption func(*initConfig) error

func WithNodeVersion(version string) InitOption {
	return func(cfg *initConfig) error {
		cfg.nodeVersion = version
		return nil
	}
}

func WithDiscovery(discovery registry.Discovery) InitOption {
	return func(cfg *initConfig) error {
		cfg.discovery = discovery
		return nil
	}
}

func WithSecret(secret string) InitOption {
	return func(cfg *initConfig) error {
		cfg.secret = secret
		return nil
	}
}

func WithProtocol(protocol string) InitOption {
	return func(cfg *initConfig) error {
		cfg.protocol = protocol
		return nil
	}
}

func WithToken(token string) InitOption {
	return func(cfg *initConfig) error {
		claims, err := rabbitMiddler.ParseClaimsFromToken(cfg.secret, token)
		if err != nil {
			return err
		}
		cfg.claim = claims
		return nil
	}
}
