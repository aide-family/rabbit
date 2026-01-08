package all

import (
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/cmd/run"
)

type Flags struct {
	*run.RunFlags

	httpTimeout    time.Duration
	grpcTimeout    time.Duration
	jobTimeout     time.Duration
	jobCoreTimeout time.Duration
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.RunFlags = run.GetRunFlags()
	c.Flags().StringVar(&f.Server.Http.Address, "http-address", f.Server.Http.Address, `Example: --http-address="0.0.0.0:8080", --http-address=":8080"`)
	c.Flags().StringVar(&f.Server.Http.Network, "http-network", f.Server.Http.Network, `Example: --http-network="tcp"`)
	c.Flags().DurationVar(&f.httpTimeout, "http-timeout", f.Server.Http.Timeout.AsDuration(), `Example: --http-timeout="10s", --http-timeout="1m", --http-timeout="1h", --http-timeout="1d"`)

	c.Flags().StringVar(&f.Server.Grpc.Address, "grpc-address", f.Server.Grpc.Address, `Example: --grpc-address="0.0.0.0:9090", --grpc-address=":9090"`)
	c.Flags().StringVar(&f.Server.Grpc.Network, "grpc-network", f.Server.Grpc.Network, `Example: --grpc-network="tcp"`)
	c.Flags().DurationVar(&f.grpcTimeout, "grpc-timeout", f.Server.Grpc.Timeout.AsDuration(), `Example: --grpc-timeout="10s", --grpc-timeout="1m", --grpc-timeout="1h", --grpc-timeout="1d"`)

	c.Flags().StringVar(&f.Server.Job.Address, "job-address", f.Server.Job.Address, `Example: --job-address="0.0.0.0:9091", --job-address=":9091"`)
	c.Flags().StringVar(&f.Server.Job.Network, "job-network", f.Server.Job.Network, `Example: --job-network="tcp"`)
	c.Flags().DurationVar(&f.jobTimeout, "job-timeout", f.Server.Job.Timeout.AsDuration(), `Example: --job-timeout="10s", --job-timeout="1m", --job-timeout="1h", --job-timeout="1d"`)

	c.Flags().Int32Var(&f.JobCore.WorkerTotal, "job-core-worker-total", f.JobCore.WorkerTotal, `Example: --job-core-worker-total=10"`)
	c.Flags().DurationVar(&f.jobCoreTimeout, "job-core-timeout", f.JobCore.Timeout.AsDuration(), `Example: --job-core-timeout="10s", --job-core-timeout="1m", --job-core-timeout="1h", --job-core-timeout="1d"`)
	c.Flags().Uint32Var(&f.JobCore.BufferSize, "job-core-buffer-size", f.JobCore.BufferSize, `Example: --job-core-buffer-size=1000"`)
}

func (f *Flags) applyToBootstrap() error {
	if err := f.ApplyToBootstrap(); err != nil {
		return err
	}
	if f.httpTimeout > 0 {
		f.Server.Http.Timeout = durationpb.New(f.httpTimeout)
	}
	if f.grpcTimeout > 0 {
		f.Server.Grpc.Timeout = durationpb.New(f.grpcTimeout)
	}
	if f.jobTimeout > 0 {
		f.Server.Job.Timeout = durationpb.New(f.jobTimeout)
	}

	if f.jobCoreTimeout > 0 {
		f.JobCore.Timeout = durationpb.New(f.jobCoreTimeout)
	}
	return nil
}
