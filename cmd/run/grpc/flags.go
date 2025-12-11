package grpc

import (
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/cmd/run"
)

type Flags struct {
	run.RunFlags
	grpcTimeout string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.RunFlags = run.GetRunFlags()
	c.Flags().StringVar(&f.Server.Grpc.Address, "grpc-address", f.Server.Grpc.Address, `Example: --grpc-address="0.0.0.0:9090", --grpc-address=":9090"`)
	c.Flags().StringVar(&f.Server.Grpc.Network, "grpc-network", f.Server.Grpc.Network, `Example: --grpc-network="tcp"`)
	c.Flags().StringVar(&f.grpcTimeout, "grpc-timeout", f.Server.Grpc.Timeout.AsDuration().String(), `Example: --grpc-timeout="10s", --grpc-timeout="1m", --grpc-timeout="1h", --grpc-timeout="1d"`)
}

func (f *Flags) applyToBootstrap() {
	if strutil.IsNotEmpty(f.grpcTimeout) {
		if timeout, err := time.ParseDuration(f.grpcTimeout); pointer.IsNil(err) {
			f.Server.Grpc.Timeout = durationpb.New(timeout)
		}
	}
}
