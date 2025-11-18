package run

import (
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/enum"
)

type Flags struct {
	cmd.GlobalFlags
	configPath             string
	clientConfigOutputPath string
	enableSwagger          bool
	enableMetrics          bool

	httpAddress string
	httpNetwork string
	httpTimeout string
	grpcAddress string
	grpcNetwork string
	grpcTimeout string
	environment string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	c.Flags().StringVarP(&f.configPath, "config", "c", "./config", "config file (default is ./config)")
	c.Flags().StringVarP(&f.clientConfigOutputPath, "client-config-output", "o", "~/.rabbit/config.yaml", "client config output file (default is ~/.rabbit/config.yaml)")
	c.Flags().BoolVarP(&f.enableSwagger, "swagger", "s", false, "enable swagger")
	c.Flags().BoolVarP(&f.enableMetrics, "metrics", "m", false, "enable metrics")
	c.Flags().StringVar(&f.httpAddress, "http-address", "0.0.0.0:8080", "http address (default is 0.0.0.0:8080)")
	c.Flags().StringVar(&f.httpNetwork, "http-network", "tcp", "http network (default is tcp)")
	c.Flags().StringVar(&f.httpTimeout, "http-timeout", "10s", "http timeout (default is 10s)")
	c.Flags().StringVar(&f.grpcAddress, "grpc-address", "0.0.0.0:9090", "grpc address (default is 0.0.0.0:9090)")
	c.Flags().StringVar(&f.grpcNetwork, "grpc-network", "tcp", "grpc network (default is tcp)")
	c.Flags().StringVar(&f.grpcTimeout, "grpc-timeout", "10s", "grpc timeout (default is 10s)")
	c.Flags().StringVarP(&f.environment, "environment", "e", "PROD", "environment (DEV, TEST, PREVIEW, PROD)")
}

func (f *Flags) applyToBootstrap(bc *conf.Bootstrap) {
	if pointer.IsNil(bc.GetServer()) {
		bc.Server = &conf.Server{
			Http: &conf.Server_HTTPServer{},
			Grpc: &conf.Server_GRPCServer{},
		}
	}
	httpConf := bc.GetServer().GetHttp()
	if pointer.IsNil(httpConf) {
		httpConf = &conf.Server_HTTPServer{}
		bc.Server.Http = httpConf
	}
	if strutil.IsEmpty(httpConf.Address) {
		httpConf.Address = f.httpAddress
	}
	if strutil.IsEmpty(httpConf.Network) {
		httpConf.Network = f.httpNetwork
	}
	if httpConf.GetTimeout().AsDuration() <= 0 {
		timeout, err := time.ParseDuration(f.httpTimeout)
		if pointer.IsNil(err) {
			httpConf.Timeout = durationpb.New(timeout)
		}
	}
	grpcConf := bc.GetServer().GetGrpc()
	if pointer.IsNil(grpcConf) {
		grpcConf = &conf.Server_GRPCServer{}
		bc.Server.Grpc = grpcConf
	}
	if strutil.IsEmpty(grpcConf.Address) {
		grpcConf.Address = f.grpcAddress
	}
	if strutil.IsEmpty(grpcConf.Network) {
		grpcConf.Network = f.grpcNetwork
	}
	if grpcConf.GetTimeout().AsDuration() <= 0 {
		timeout, err := time.ParseDuration(f.grpcTimeout)
		if pointer.IsNil(err) {
			grpcConf.Timeout = durationpb.New(timeout)
		}
	}
	if bc.Environment.IsUnknown() {
		env := enum.Environment_PROD
		if strutil.IsNotEmpty(f.environment) {
			env = enum.Environment(enum.Environment_value[f.environment])
		}
		bc.Environment = env
	}
}
