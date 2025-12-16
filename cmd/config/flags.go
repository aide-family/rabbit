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
	c.Flags().StringVarP(&f.name, "name", "N", "config.yaml", "output file name (default is config.yaml)")
	c.Flags().BoolVarP(&f.force, "force", "f", false, "overwrite existing file if it exists (default is to rename with timestamp)")
	c.Flags().BoolVar(&f.isClient, "client", false, "generate client config file (default is server config file)")
}

var clientConfig = &config.ClientConfig{
	RegistryType: config.RegistryType_UNKNOWN,
	Cluster: &config.ClusterConfig{
		Name:      "rabbit",
		Endpoints: "localhost:8080",
		Timeout:   durationpb.New(10 * time.Second),
	},
	JwtToken: "Bearer <jwt-token>",
}
