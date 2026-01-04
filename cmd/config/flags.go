package config

import (
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/pkg/config"
)

type Flags struct {
	*cmd.GlobalFlags
	path     string
	name     string
	force    bool
	isClient bool
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.GlobalFlags = cmd.GetGlobalFlags()
	c.Flags().StringVarP(&f.path, "path", "p", ".", "output path for the config file (default is current directory)")
	c.Flags().StringVarP(&f.name, "name", "N", "rabbit_config.yaml", "output file name (default is config.yaml)")
	c.Flags().BoolVarP(&f.force, "force", "f", false, "overwrite existing file if it exists (default is to rename with timestamp)")
	c.Flags().BoolVar(&f.isClient, "client", false, "generate client config file (default is server config file)")
}

var clientConfig = &config.ClientConfig{
	RegistryType: config.RegistryType_UNKNOWN,
	Cluster: &config.ClusterConfig{
		Name:      "moon.rabbit",
		Endpoints: "localhost:10080",
		Timeout:   durationpb.New(10 * time.Second),
		Protocol:  config.ClusterConfig_GRPC,
	},
	JwtToken: "Bearer <jwt-token>",
	Etcd: &config.ETCDConfig{
		Endpoints: "localhost:2379",
		Username:  "",
		Password:  "",
	},
	Kubernetes: &config.KubernetesConfig{
		KubeConfig: "~/.kube/config",
	},
	Namespace: "moon",
}
