package job

import (
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/cmd/run"
)

type Flags struct {
	*run.RunFlags

	jobTimeout     string
	jobCoreTimeout string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.RunFlags = run.GetRunFlags()
	c.Flags().StringVar(&f.Server.Job.Address, "job-address", f.Server.Job.Address, `Example: --job-address="0.0.0.0:9091", --job-address=":9091"`)
	c.Flags().StringVar(&f.Server.Job.Network, "job-network", f.Server.Job.Network, `Example: --job-network="grpc", --job-network="http"`)
	c.Flags().StringVar(&f.jobTimeout, "job-timeout", f.Server.Job.Timeout.AsDuration().String(), `Example: --job-timeout="10s", --job-timeout="1m", --job-timeout="1h", --job-timeout="1d"`)
	c.Flags().Int32Var(&f.JobCore.WorkerTotal, "job-core-worker-total", f.JobCore.WorkerTotal, `Example: --job-core-worker-total=10"`)
	c.Flags().StringVar(&f.jobCoreTimeout, "job-core-timeout", f.JobCore.Timeout.AsDuration().String(), `Example: --job-core-timeout="10s", --job-core-timeout="1m", --job-core-timeout="1h", --job-core-timeout="1d"`)
	c.Flags().Uint32Var(&f.JobCore.BufferSize, "job-core-buffer-size", f.JobCore.BufferSize, `Example: --job-core-buffer-size=1000"`)
}

func (f *Flags) applyToBootstrap() {
	f.ApplyToBootstrap()
	if strutil.IsNotEmpty(f.jobTimeout) {
		if timeout, err := time.ParseDuration(f.jobTimeout); pointer.IsNil(err) {
			f.Server.Job.Timeout = durationpb.New(timeout)
		}
	}

	if strutil.IsNotEmpty(f.jobCoreTimeout) {
		if timeout, err := time.ParseDuration(f.jobCoreTimeout); pointer.IsNil(err) {
			f.JobCore.Timeout = durationpb.New(timeout)
		}
	}
}
