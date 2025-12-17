package job

import (
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/cmd/run"
)

type Flags struct {
	*run.RunFlags

	jobTimeout     time.Duration
	jobCoreTimeout time.Duration
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.RunFlags = run.GetRunFlags()
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
	if f.jobTimeout > 0 {
		f.Server.Job.Timeout = durationpb.New(f.jobTimeout)
	}

	if f.jobCoreTimeout > 0 {
		f.JobCore.Timeout = durationpb.New(f.jobCoreTimeout)
	}
	return nil
}
