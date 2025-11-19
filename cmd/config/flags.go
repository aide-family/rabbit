package config

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

type Flags struct {
	cmd.GlobalFlags
	path  string
	name  string
	force bool
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.GlobalFlags = cmd.GetGlobalFlags()
	c.Flags().StringVarP(&f.path, "path", "p", ".", "output path for the config file (default is current directory)")
	c.Flags().StringVar(&f.name, "name", "server.yaml", "output file name (default is server.yaml)")
	c.Flags().BoolVarP(&f.force, "force", "f", false, "overwrite existing file if it exists (default is to rename with timestamp)")
}
