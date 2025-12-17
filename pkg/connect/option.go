package connect

import (
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/merr"
)

func init() {
	selector.SetGlobalSelector(p2c.NewBuilder())
}

var (
	ProtocolHTTP = config.ClusterConfig_HTTP.String()
	ProtocolGRPC = config.ClusterConfig_GRPC.String()
)

type InitConfig interface {
	GetName() string
	GetEndpoint() string
	GetTimeout() *durationpb.Duration
	GetProtocol() string
}

type DefaultConfig struct {
	name     string
	endpoint string
	timeout  time.Duration
	protocol string
}

func NewDefaultConfig(name, endpoint string, timeout time.Duration, protocol string) InitConfig {
	return &DefaultConfig{
		name:     name,
		endpoint: endpoint,
		timeout:  timeout,
		protocol: protocol,
	}
}

func (c *DefaultConfig) GetName() string {
	return c.name
}

func (c *DefaultConfig) GetEndpoint() string {
	return c.endpoint
}

func (c *DefaultConfig) GetTimeout() *durationpb.Duration {
	return durationpb.New(c.timeout)
}

func (c *DefaultConfig) GetProtocol() string {
	return c.protocol
}

type initConfig struct {
	name        string
	endpoint    string
	protocol    string
	timeout     time.Duration
	nodeVersion string
	discovery   registry.Discovery
	nodeFilters []NodeFilter
}

func NewInitConfig(config InitConfig, opts ...InitOption) (*initConfig, error) {
	cfg := &initConfig{
		name:        config.GetName(),
		endpoint:    config.GetEndpoint(),
		protocol:    config.GetProtocol(),
		timeout:     config.GetTimeout().AsDuration(),
		nodeFilters: []NodeFilter{},
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

func WithNodeFilter(filter func(node selector.Node) bool) InitOption {
	return func(cfg *initConfig) error {
		if pointer.IsNotNil(cfg.nodeFilters) {
			return merr.ErrorInternalServer("node filter already set")
		}
		cfg.nodeFilters = append(cfg.nodeFilters, filter)
		return nil
	}
}
