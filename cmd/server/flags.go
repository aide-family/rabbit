package server

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

	httpAddress string
	httpNetwork string
	httpTimeout string
	grpcAddress string
	grpcNetwork string
	grpcTimeout string
	environment string

	enableHTTP bool
	enableGRPC bool
	enableJob  bool

	flagSet *pflag.FlagSet
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.flagSet = c.Flags()
	c.Flags().StringVarP(&f.configPath, "config", "c", "./config", "config file (default is ./config)")
	c.Flags().StringVarP(&f.clientConfigOutputPath, "client-config-output", "o", "~/.rabbit/config.yaml", "client config output file (default is ~/.rabbit/config.yaml)")

	c.Flags().StringVar(&f.httpAddress, "http-address", "0.0.0.0:8080", "http address (default is 0.0.0.0:8080)")
	c.Flags().StringVar(&f.httpNetwork, "http-network", "tcp", "http network (default is tcp)")
	c.Flags().StringVar(&f.httpTimeout, "http-timeout", "10s", "http timeout (default is 10s)")
	c.Flags().StringVar(&f.grpcAddress, "grpc-address", "0.0.0.0:9090", "grpc address (default is 0.0.0.0:9090)")
	c.Flags().StringVar(&f.grpcNetwork, "grpc-network", "tcp", "grpc network (default is tcp)")
	c.Flags().StringVar(&f.grpcTimeout, "grpc-timeout", "10s", "grpc timeout (default is 10s)")
	c.Flags().StringVarP(&f.environment, "environment", "e", "PROD", "environment (DEV, TEST, PREVIEW, PROD)")

	c.Flags().BoolVar(&f.enableHTTP, "http", false, "enable HTTP server")
	c.Flags().BoolVar(&f.enableGRPC, "grpc", false, "enable gRPC server")
	c.Flags().BoolVar(&f.enableJob, "job", false, "enable job workers")
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

	isFlagChanged := func(name string) bool {
		if f.flagSet == nil {
			return false
		}
		if flag := f.flagSet.Lookup(name); flag != nil {
			return flag.Changed
		}
		return false
	}
	setBoolString := func(target *string, value bool) {
		if value {
			*target = "true"
		} else {
			*target = "false"
		}
	}

	httpChanged := isFlagChanged("http")
	grpcChanged := isFlagChanged("grpc")
	jobChanged := isFlagChanged("job")
	componentChanged := httpChanged || grpcChanged || jobChanged

	if componentChanged {
		if httpChanged {
			setBoolString(&bc.Server.EnableHttp, f.enableHTTP)
		} else {
			setBoolString(&bc.Server.EnableHttp, false)
		}
		if grpcChanged {
			setBoolString(&bc.Server.EnableGrpc, f.enableGRPC)
		} else {
			setBoolString(&bc.Server.EnableGrpc, false)
		}
		if jobChanged {
			setBoolString(&bc.Server.EnableJob, f.enableJob)
		} else {
			setBoolString(&bc.Server.EnableJob, false)
		}
		return
	}

	if strutil.IsEmpty(bc.Server.EnableHttp) {
		setBoolString(&bc.Server.EnableHttp, true)
	}
	if strutil.IsEmpty(bc.Server.EnableGrpc) {
		setBoolString(&bc.Server.EnableGrpc, true)
	}
	if strutil.IsEmpty(bc.Server.EnableJob) {
		setBoolString(&bc.Server.EnableJob, false)
	}
}
