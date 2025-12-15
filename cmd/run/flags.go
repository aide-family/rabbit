package run

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/enum"
)

type RunFlags struct {
	*conf.Bootstrap
	*cmd.GlobalFlags

	metadata           []string
	useRandomID        bool
	configPaths        []string
	dataSourcePaths    []string
	environment        string
	jwtExpire          string
	registryType       string
	enableClientConfig bool
}

var runFlags RunFlags

func (f *RunFlags) addFlags(c *cobra.Command, bc *conf.Bootstrap) {
	f.GlobalFlags = cmd.GetGlobalFlags()
	f.Bootstrap = bc

	c.PersistentFlags().StringSliceVarP(&f.configPaths, "config", "c", []string{}, `Example: -c=./config1/ -c=./config2/`)
	enableClientConfig, _ := strconv.ParseBool(f.EnableClientConfig)
	c.PersistentFlags().BoolVar(&f.enableClientConfig, "enable-client-config", enableClientConfig, `Example: --enable-client-config`)

	c.PersistentFlags().StringVar(&f.Server.Name, "server-name", f.Server.Name, `Example: --server-name="rabbit"`)
	useRandomID, _ := strconv.ParseBool(f.Server.UseRandomID)
	c.PersistentFlags().BoolVar(&f.useRandomID, "use-random-node-id", useRandomID, `Example: --use-random-node-id`)
	metadataStr := make([]string, 0, len(f.Server.Metadata))
	for key, value := range f.Server.Metadata {
		metadataStr = append(metadataStr, fmt.Sprintf("%s=%s", key, value))
	}
	c.PersistentFlags().StringSliceVar(&f.metadata, "server-metadata", metadataStr, `Example: --server-metadata="tag=rabbit" --server-metadata="email=aidecloud@163.com"`)
	c.PersistentFlags().StringVar(&f.environment, "environment", f.Environment.String(), `Example: --environment="DEV", --environment="TEST", --environment="PREVIEW", --environment="PROD"`)
	c.PersistentFlags().StringVar(&f.Jwt.Secret, "jwt-secret", f.Jwt.Secret, `Example: --jwt-secret="xxx"`)
	c.PersistentFlags().StringVar(&f.jwtExpire, "jwt-expire", f.Jwt.Expire.AsDuration().String(), `Example: --jwt-expire="10s", --jwt-expire="1m", --jwt-expire="1h", --jwt-expire="1d"`)
	c.PersistentFlags().StringVar(&f.Jwt.Issuer, "jwt-issuer", f.Jwt.Issuer, `Example: --jwt-issuer="xxx"`)
	c.PersistentFlags().StringVar(&f.Main.Username, "main-username", f.Main.Username, `Example: --main-username="root"`)
	c.PersistentFlags().StringVar(&f.Main.Password, "main-password", f.Main.Password, `Example: --main-password="123456"`)
	c.PersistentFlags().StringVar(&f.Main.Host, "main-host", f.Main.Host, `Example: --main-host="localhost"`)
	c.PersistentFlags().Int32Var(&f.Main.Port, "main-port", f.Main.Port, `Example: --main-port=3306"`)
	c.PersistentFlags().StringVar(&f.Main.Database, "main-database", f.Main.Database, `Example: --main-database="rabbit"`)
	c.PersistentFlags().StringVar(&f.Main.Debug, "main-debug", f.Main.Debug, `Example: --main-debug="false"`)
	c.PersistentFlags().StringVar(&f.Main.UseSystemLogger, "main-use-system-logger", f.Main.UseSystemLogger, `Example: --main-use-system-logger="true"`)
	c.PersistentFlags().StringVar(&f.registryType, "registry-type", f.RegistryType.String(), `Example: --registry-type="ETCD"`)
	c.PersistentFlags().StringVar(&f.Etcd.Endpoints, "etcd-endpoints", f.Etcd.Endpoints, `Example: --etcd-endpoints="127.0.0.1:2379"`)
	c.PersistentFlags().StringVar(&f.Etcd.Username, "etcd-username", f.Etcd.Username, `Example: --etcd-username="root"`)
	c.PersistentFlags().StringVar(&f.Etcd.Password, "etcd-password", f.Etcd.Password, `Example: --etcd-password="123456"`)
	c.PersistentFlags().StringVar(&f.Kubernetes.KubeConfig, "kubernetes-kubeconfig", f.Kubernetes.KubeConfig, `Example: --kubernetes-kubeconfig="~/.kube/config"`)
	c.PersistentFlags().StringVar(&f.UseDatabase, "use-database", f.UseDatabase, `Example: --use-database="true"`)
	c.PersistentFlags().StringSliceVar(&f.dataSourcePaths, "datasource-paths", strutil.SplitSkipEmpty(f.DataSourcePaths, ","), `Example: --datasource-paths="./datasource" --datasource-paths="./config,./datasource"`)
	c.PersistentFlags().StringVar(&f.MessageLogPath, "message-log-path", f.MessageLogPath, `Example: --message-log-path="./messages/"`)
}

func (f *RunFlags) ApplyToBootstrap() {
	if strutil.IsEmpty(f.Server.Name) {
		f.Server.Name = f.Name
	}
	if strutil.IsEmpty(f.Server.Namespace) {
		f.Server.Namespace = f.Namespace
	}
	f.EnableClientConfig = strconv.FormatBool(f.enableClientConfig)

	metadata := f.Server.Metadata
	if pointer.IsNil(metadata) {
		metadata = make(map[string]string)
	}

	metadata["repository"] = f.Repo
	metadata["author"] = f.Author
	metadata["email"] = f.Email
	metadata["built"] = f.Built

	for _, m := range f.metadata {
		parts := strings.SplitN(m, "=", 2)
		if len(parts) == 2 {
			metadata[parts[0]] = parts[1]
		}
	}

	f.Server.Metadata = metadata
	f.Server.UseRandomID = strconv.FormatBool(f.useRandomID)

	if strutil.IsNotEmpty(f.environment) {
		f.Environment = enum.Environment(enum.Environment_value[f.environment])
	}

	if strutil.IsNotEmpty(f.jwtExpire) {
		if expire, err := time.ParseDuration(f.jwtExpire); pointer.IsNil(err) {
			f.Jwt.Expire = durationpb.New(expire)
		}
	}

	if strutil.IsNotEmpty(f.registryType) {
		f.RegistryType = config.RegistryType(config.RegistryType_value[f.registryType])
	}

	if len(f.configPaths) > 0 {
		var bc conf.Bootstrap
		sourceOpts := make([]kconfig.Source, 0, len(f.configPaths))
		sourceOpts = append(sourceOpts, env.NewSource())
		for _, configPath := range f.configPaths {
			if strutil.IsNotEmpty(configPath) {
				sourceOpts = append(sourceOpts, file.NewSource(load.ExpandHomeDir(strings.TrimSpace(configPath))))
			}
		}
		if len(sourceOpts) > 0 {
			if err := conf.Load(&bc, sourceOpts...); err != nil {
				klog.Errorw("msg", "load config failed", "error", err)
				return
			}
			f.Bootstrap = &bc
		}
	}
	if len(f.dataSourcePaths) > 0 {
		f.DataSourcePaths = strings.Join(f.configPaths, ",")
	}
}

func GetRunFlags() *RunFlags {
	return &runFlags
}

type WireApp func(serviceName string, bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, func(), error)

func StartServer(serviceName string, wireApp WireApp) {
	serverConf := runFlags.GetServer()
	envOpts := []hello.Option{
		hello.WithVersion(runFlags.Version),
		hello.WithID(runFlags.Hostname),
		hello.WithEnv(runFlags.Environment.String()),
		hello.WithMetadata(serverConf.GetMetadata()),
	}
	if strings.EqualFold(serverConf.GetUseRandomID(), "true") {
		envOpts = append(envOpts, hello.WithID(strutil.RandomID()))
	}
	hello.SetEnvWithOption(envOpts...)
	helper := klog.NewHelper(klog.With(klog.GetLogger(),
		"service.name", serviceName,
		"service.id", hello.ID(),
		"caller", klog.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID()),
	)

	app, cleanup, err := wireApp(serviceName, runFlags.Bootstrap, helper)
	if err != nil {
		klog.Errorw("msg", "wireApp failed", "error", err)
		return
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		klog.Errorw("msg", "app run failed", "error", err)
		return
	}
}
