package all

import (
	"strconv"
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/cmd/run"
)

type Flags struct {
	*run.RunFlags

	httpTimeout            string
	grpcTimeout            string
	jobTimeout             string
	jobCoreTimeout         string
	enableSwagger          bool
	enableSwaggerBasicAuth bool
	enableMetrics          bool
	enableMetricsBasicAuth bool
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.RunFlags = run.GetRunFlags()
	c.Flags().StringVar(&f.Server.Http.Address, "http-address", f.Server.Http.Address, `Example: --http-address="0.0.0.0:8080", --http-address=":8080"`)
	c.Flags().StringVar(&f.Server.Http.Network, "http-network", f.Server.Http.Network, `Example: --http-network="tcp"`)
	c.Flags().StringVar(&f.httpTimeout, "http-timeout", f.Server.Http.Timeout.AsDuration().String(), `Example: --http-timeout="10s", --http-timeout="1m", --http-timeout="1h", --http-timeout="1d"`)
	enableSwagger, _ := strconv.ParseBool(f.SwaggerBasicAuth.Enabled)
	c.Flags().BoolVar(&f.enableSwagger, "enable-swagger", enableSwagger, `Example: --enable-swagger`)
	enableSwaggerBasicAuth, _ := strconv.ParseBool(f.SwaggerBasicAuth.Enabled)
	c.Flags().BoolVar(&f.enableSwaggerBasicAuth, "enable-swagger-basic-auth", enableSwaggerBasicAuth, `Example: --enable-swagger-basic-auth`)
	c.Flags().StringVar(&f.SwaggerBasicAuth.Username, "swagger-basic-auth-username", f.SwaggerBasicAuth.Username, `Example: --swagger-basic-auth-username="username"`)
	c.Flags().StringVar(&f.SwaggerBasicAuth.Password, "swagger-basic-auth-password", f.SwaggerBasicAuth.Password, `Example: --swagger-basic-auth-password="password"`)
	enableMetrics, _ := strconv.ParseBool(f.MetricsBasicAuth.Enabled)
	c.Flags().BoolVar(&f.enableMetrics, "enable-metrics", enableMetrics, `Example: --enable-metrics`)
	enableMetricsBasicAuth, _ := strconv.ParseBool(f.MetricsBasicAuth.Enabled)
	c.Flags().BoolVar(&f.enableMetricsBasicAuth, "enable-metrics-basic-auth", enableMetricsBasicAuth, `Example: --enable-metrics-basic-auth`)
	c.Flags().StringVar(&f.MetricsBasicAuth.Username, "metrics-basic-auth-username", f.MetricsBasicAuth.Username, `Example: --metrics-basic-auth-username="username"`)
	c.Flags().StringVar(&f.MetricsBasicAuth.Password, "metrics-basic-auth-password", f.MetricsBasicAuth.Password, `Example: --metrics-basic-auth-password="password"`)

	c.Flags().StringVar(&f.Server.Grpc.Address, "grpc-address", f.Server.Grpc.Address, `Example: --grpc-address="0.0.0.0:9090", --grpc-address=":9090"`)
	c.Flags().StringVar(&f.Server.Grpc.Network, "grpc-network", f.Server.Grpc.Network, `Example: --grpc-network="tcp"`)
	c.Flags().StringVar(&f.grpcTimeout, "grpc-timeout", f.Server.Grpc.Timeout.AsDuration().String(), `Example: --grpc-timeout="10s", --grpc-timeout="1m", --grpc-timeout="1h", --grpc-timeout="1d"`)

	c.Flags().StringVar(&f.Server.Job.Address, "job-address", f.Server.Job.Address, `Example: --job-address="0.0.0.0:9091", --job-address=":9091"`)
	c.Flags().StringVar(&f.Server.Job.Network, "job-network", f.Server.Job.Network, `Example: --job-network="tcp"`)
	c.Flags().StringVar(&f.jobTimeout, "job-timeout", f.Server.Job.Timeout.AsDuration().String(), `Example: --job-timeout="10s", --job-timeout="1m", --job-timeout="1h", --job-timeout="1d"`)

	c.Flags().Int32Var(&f.JobCore.WorkerTotal, "job-core-worker-total", f.JobCore.WorkerTotal, `Example: --job-core-worker-total=10"`)
	c.Flags().StringVar(&f.jobCoreTimeout, "job-core-timeout", f.JobCore.Timeout.AsDuration().String(), `Example: --job-core-timeout="10s", --job-core-timeout="1m", --job-core-timeout="1h", --job-core-timeout="1d"`)
	c.Flags().Uint32Var(&f.JobCore.BufferSize, "job-core-buffer-size", f.JobCore.BufferSize, `Example: --job-core-buffer-size=1000"`)
}

func (f *Flags) applyToBootstrap() {
	f.ApplyToBootstrap()
	if strutil.IsNotEmpty(f.httpTimeout) {
		if timeout, err := time.ParseDuration(f.httpTimeout); pointer.IsNil(err) {
			f.Server.Http.Timeout = durationpb.New(timeout)
		}
	}
	if strutil.IsNotEmpty(f.grpcTimeout) {
		if timeout, err := time.ParseDuration(f.grpcTimeout); pointer.IsNil(err) {
			f.Server.Grpc.Timeout = durationpb.New(timeout)
		}
	}
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
