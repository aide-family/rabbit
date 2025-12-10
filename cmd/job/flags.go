package job

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
	configPaths []string
	useEnv      bool

	*conf.Bootstrap
	environment    string
	jobTimeout     string
	jwtExpire      string
	jobCoreTimeout string
	registryType   string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command, bc *conf.Bootstrap) {
	f.Bootstrap = bc
	c.Flags().StringSliceVarP(&f.configPaths, "config", "c", []string{}, `Example: -c=./config1/ -c=./config2/`)
	c.Flags().BoolVar(&f.useEnv, "use-env", false, `Example: --use-env or --use-env=true`)

	c.Flags().StringVar(&f.environment, "environment", f.Environment.String(), `Example: --environment="DEV", --environment="TEST", --environment="PREVIEW", --environment="PROD"`)
	c.Flags().StringVar(&f.Server.Job.Address, "job-address", f.Server.Job.Address, `Example: --job-address="0.0.0.0:9091", --job-address=":9091"`)
	c.Flags().StringVar(&f.Server.Job.Network, "job-network", f.Server.Job.Network, `Example: --job-network="tcp"`)
	c.Flags().StringVar(&f.jobTimeout, "job-timeout", f.Server.Job.Timeout.AsDuration().String(), `Example: --job-timeout="10s", --job-timeout="1m", --job-timeout="1h", --job-timeout="1d"`)
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
	c.Flags().Int32Var(&f.JobCore.WorkerTotal, "job-core-worker-total", f.JobCore.WorkerTotal, `Example: --job-core-worker-total=10"`)
	c.Flags().StringVar(&f.jobCoreTimeout, "job-core-timeout", f.JobCore.Timeout.AsDuration().String(), `Example: --job-core-timeout="10s", --job-core-timeout="1m", --job-core-timeout="1h", --job-core-timeout="1d"`)
	c.Flags().Uint32Var(&f.JobCore.BufferSize, "job-core-buffer-size", f.JobCore.BufferSize, `Example: --job-core-buffer-size=1000"`)
	c.Flags().StringVar(&f.registryType, "registry-type", f.RegistryType.String(), `Example: --registry-type="etcd"`)
	c.Flags().StringVar(&f.Etcd.Endpoints, "etcd-endpoints", f.Etcd.Endpoints, `Example: --etcd-endpoints="127.0.0.1:2379"`)
	c.Flags().StringVar(&f.Etcd.Username, "etcd-username", f.Etcd.Username, `Example: --etcd-username="root"`)
	c.Flags().StringVar(&f.Etcd.Password, "etcd-password", f.Etcd.Password, `Example: --etcd-password="123456"`)
	c.Flags().StringVar(&f.Kubernetes.Namespace, "kubernetes-namespace", f.Kubernetes.Namespace, `Example: --kubernetes-namespace="moon"`)
	c.Flags().StringVar(&f.Kubernetes.KubeConfig, "kubernetes-kubeconfig", f.Kubernetes.KubeConfig, `Example: --kubernetes-kubeconfig="~/.kube/config"`)
	c.Flags().StringVar(&f.UseDatabase, "use-database", f.UseDatabase, `Example: --use-database="true"`)
	c.Flags().StringVar(&f.ConfigPaths, "config-paths", f.ConfigPaths, `Example: --config-paths="./datasource" --config-paths="./config,./datasource"`)
	c.Flags().StringVar(&f.MessageLogPath, "message-log-path", f.MessageLogPath, `Example: --message-log-path="./messages/"`)
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
	if strutil.IsNotEmpty(f.jobTimeout) {
		if timeout, err := time.ParseDuration(f.jobTimeout); pointer.IsNil(err) {
			f.Server.Job.Timeout = durationpb.New(timeout)
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
	if strutil.IsNotEmpty(f.jobCoreTimeout) {
		if timeout, err := time.ParseDuration(f.jobCoreTimeout); pointer.IsNil(err) {
			f.JobCore.Timeout = durationpb.New(timeout)
		}
	}

	if strutil.IsNotEmpty(f.registryType) {
		f.RegistryType = config.RegistryType(config.RegistryType_value[f.registryType])
	}
}
