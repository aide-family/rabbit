package grpc

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
	configPath string

	grpcAddress string
	grpcNetwork string
	grpcTimeout string
	environment string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	c.Flags().StringVarP(&f.configPath, "config", "c", "./config", "config file (default is ./config)")

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

	// this command is gRPC-only
	bc.Server.EnableHttp = "false"
	bc.Server.EnableGrpc = "true"
	bc.Server.EnableJob = "false"
}
