package run

import (
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/enum"
)

type Flags struct {
	cmd.GlobalFlags
	configPath             string
	clientConfigOutputPath string

	*conf.Bootstrap
	environment     string
	httpTimeout     string
	grpcTimeout     string
	jwtExpire       string
	eventBusTimeout string
	registryType    string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command, bc *conf.Bootstrap) {
	f.Bootstrap = bc
	c.Flags().StringVarP(&f.configPath, "config", "c", "", `Example: -c=./config/, --config="./config/"`)
	c.Flags().StringVarP(&f.clientConfigOutputPath, "client-config-output", "o", "~/.rabbit/", `Example: -o=./client/, --client-config-output="~/.rabbit/"`)

	c.Flags().StringVar(&f.environment, "environment", f.Environment.String(), `Example: --environment="DEV", --environment="TEST", --environment="PREVIEW", --environment="PROD"`)
	c.Flags().StringVar(&f.Server.Http.Address, "http-address", f.Server.Http.Address, `Example: --http-address="0.0.0.0:8080", --http-address=":8080"`)
	c.Flags().StringVar(&f.Server.Http.Network, "http-network", f.Server.Http.Network, `Example: --http-network="tcp"`)
	c.Flags().StringVar(&f.httpTimeout, "http-timeout", f.Server.Http.Timeout.AsDuration().String(), `Example: --http-timeout="10s", --http-timeout="1m", --http-timeout="1h", --http-timeout="1d"`)
	c.Flags().StringVar(&f.Server.Grpc.Address, "grpc-address", f.Server.Grpc.Address, `Example: --grpc-address="0.0.0.0:9090", --grpc-address=":9090"`)
	c.Flags().StringVar(&f.Server.Grpc.Network, "grpc-network", f.Server.Grpc.Network, `Example: --grpc-network="tcp"`)
	c.Flags().StringVar(&f.grpcTimeout, "grpc-timeout", f.Server.Grpc.Timeout.AsDuration().String(), `Example: --grpc-timeout="10s", --grpc-timeout="1m", --grpc-timeout="1h", --grpc-timeout="1d"`)
	c.Flags().StringVar(&f.Jwt.Secret, "jwt-secret", f.Jwt.Secret, `Example: --jwt-secret="xxx"`)
	c.Flags().StringVar(&f.jwtExpire, "jwt-expire", f.Jwt.Expire.AsDuration().String(), `Example: --jwt-expire="10s", --jwt-expire="1m", --jwt-expire="1h", --jwt-expire="1d"`)
	c.Flags().StringVar(&f.Jwt.Issuer, "jwt-issuer", f.Jwt.Issuer, `Example: --jwt-issuer="xxx"`)
	c.Flags().StringVar(&f.Main.Username, "main-username", f.Main.Username, `Example: --main-username="root"`)
	c.Flags().StringVar(&f.Main.Password, "main-password", f.Main.Password, `Example: --main-password="123456"`)
	c.Flags().StringVar(&f.Main.Host, "main-host", f.Main.Host, `Example: --main-host="localhost"`)
	c.Flags().Int32Var(&f.Main.Port, "main-port", f.Main.Port, `Example: --main-port=3306"`)
	c.Flags().StringVar(&f.Main.Database, "main-database", f.Main.Database, `Example: --main-database="rabbit"`)
	c.Flags().StringVar(&f.Main.Debug, "main-debug", f.Main.Debug, `Example: --main-debug="false"`)
	c.Flags().StringVar(&f.Main.UseSystemLogger, "main-use-system-logger", f.Main.UseSystemLogger, `Example: --main-use-system-logger="true"`)
	c.Flags().Int32Var(&f.EventBus.WorkerCount, "event-bus-worker-count", f.EventBus.WorkerCount, `Example: --event-bus-worker-count=1"`)
	c.Flags().StringVar(&f.eventBusTimeout, "event-bus-timeout", f.EventBus.Timeout.AsDuration().String(), `Example: --event-bus-timeout="10s", --event-bus-timeout="1m", --event-bus-timeout="1h", --event-bus-timeout="1d"`)
	c.Flags().StringVar(&f.registryType, "registry-type", f.RegistryType.String(), `Example: --registry-type="etcd"`)
	c.Flags().StringVar(&f.Etcd.Endpoints, "etcd-endpoints", f.Etcd.Endpoints, `Example: --etcd-endpoints="127.0.0.1:2379"`)
	c.Flags().StringVar(&f.Etcd.Username, "etcd-username", f.Etcd.Username, `Example: --etcd-username="root"`)
	c.Flags().StringVar(&f.Etcd.Password, "etcd-password", f.Etcd.Password, `Example: --etcd-password="123456"`)
	c.Flags().StringVar(&f.Kubernetes.Namespace, "kubernetes-namespace", f.Kubernetes.Namespace, `Example: --kubernetes-namespace="moon"`)
	c.Flags().StringVar(&f.Kubernetes.KubeConfig, "kubernetes-kubeconfig", f.Kubernetes.KubeConfig, `Example: --kubernetes-kubeconfig="~/.kube/config"`)
	c.Flags().StringVar(&f.SwaggerBasicAuth.Username, "swagger-basic-auth-username", f.SwaggerBasicAuth.Username, `Example: --swagger-basic-auth-username="root"`)
	c.Flags().StringVar(&f.SwaggerBasicAuth.Password, "swagger-basic-auth-password", f.SwaggerBasicAuth.Password, `Example: --swagger-basic-auth-password="123456"`)
	c.Flags().StringVar(&f.MetricsBasicAuth.Username, "metrics-basic-auth-username", f.MetricsBasicAuth.Username, `Example: --metrics-basic-auth-username="root"`)
	c.Flags().StringVar(&f.MetricsBasicAuth.Password, "metrics-basic-auth-password", f.MetricsBasicAuth.Password, `Example: --metrics-basic-auth-password="123456"`)
	c.Flags().StringVar(&f.EnableClientConfig, "enable-client-config", f.EnableClientConfig, `Example: --enable-client-config="true"`)
	c.Flags().StringVar(&f.EnableSwagger, "enable-swagger", f.EnableSwagger, `Example: --enable-swagger="true"`)
	c.Flags().StringVar(&f.EnableMetrics, "enable-metrics", f.EnableMetrics, `Example: --enable-metrics="true"`)
	c.Flags().StringVar(&f.UseDatabase, "use-database", f.UseDatabase, `Example: --use-database="true"`)
	c.Flags().StringSliceVar(&f.ConfigPaths, "config-paths", f.ConfigPaths, `Example: --config-paths="./datasource"`)
}

func (f *Flags) applyToBootstrap() {
	metadata := f.Server.Metadata
	if pointer.IsNil(metadata) {
		metadata = make(map[string]string)
	}
	metadata["repository"] = f.Repo
	metadata["author"] = f.Author
	metadata["email"] = f.Email
	f.Server.Metadata = metadata
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

	if strutil.IsNotEmpty(f.environment) {
		f.Environment = enum.Environment(enum.Environment_value[f.environment])
	}
	if strutil.IsNotEmpty(f.jwtExpire) {
		if expire, err := time.ParseDuration(f.jwtExpire); pointer.IsNil(err) {
			f.Jwt.Expire = durationpb.New(expire)
		}
	}
	if strutil.IsNotEmpty(f.eventBusTimeout) {
		if timeout, err := time.ParseDuration(f.eventBusTimeout); pointer.IsNil(err) {
			f.EventBus.Timeout = durationpb.New(timeout)
		}
	}

	if strutil.IsNotEmpty(f.registryType) {
		f.RegistryType = config.RegistryType(config.RegistryType_value[f.registryType])
	}
}
