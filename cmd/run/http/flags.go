package http

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/cmd/run"
)

type Flags struct {
	*run.RunFlags

	httpTimeout            time.Duration
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
	c.Flags().DurationVar(&f.httpTimeout, "http-timeout", f.Server.Http.Timeout.AsDuration(), `Example: --http-timeout="10s", --http-timeout="1m", --http-timeout="1h", --http-timeout="1d"`)

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
}

func (f *Flags) applyToBootstrap() error {
	if err := f.ApplyToBootstrap(); err != nil {
		return err
	}
	if f.httpTimeout > 0 {
		f.Server.Http.Timeout = durationpb.New(f.httpTimeout)
	}

	f.EnableSwagger = strconv.FormatBool(f.enableSwagger)
	f.SwaggerBasicAuth.Enabled = strconv.FormatBool(f.enableSwaggerBasicAuth)
	f.EnableMetrics = strconv.FormatBool(f.enableMetrics)
	f.MetricsBasicAuth.Enabled = strconv.FormatBool(f.enableMetricsBasicAuth)
	return nil
}
