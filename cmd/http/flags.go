package http

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

	httpAddress string
	httpNetwork string
	httpTimeout string
	environment string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	c.Flags().StringVarP(&f.configPath, "config", "c", "./config", "config file (default is ./config)")

	c.Flags().StringVar(&f.httpAddress, "http-address", "0.0.0.0:8080", "http address (default is 0.0.0.0:8080)")
	c.Flags().StringVar(&f.httpNetwork, "http-network", "tcp", "http network (default is tcp)")
	c.Flags().StringVar(&f.httpTimeout, "http-timeout", "10s", "http timeout (default is 10s)")
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

	if bc.Environment.IsUnknown() {
		env := enum.Environment_PROD
		if strutil.IsNotEmpty(f.environment) {
			env = enum.Environment(enum.Environment_value[f.environment])
		}
		bc.Environment = env
	}

	// this command is HTTP-only
	bc.Server.EnableHttp = "true"
	bc.Server.EnableGrpc = "false"
	bc.Server.EnableJob = "false"
}
