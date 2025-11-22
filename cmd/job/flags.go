package job

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/enum"
)

type Flags struct {
	cmd.GlobalFlags
	configPath  string
	environment string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	c.Flags().StringVarP(&f.configPath, "config", "c", "./config", "config file (default is ./config)")
	c.Flags().StringVarP(&f.environment, "environment", "e", "PROD", "environment (DEV, TEST, PREVIEW, PROD)")
}

func (f *Flags) applyToBootstrap(bc *conf.Bootstrap) {
	if bc.Environment.IsUnknown() {
		env := enum.Environment_PROD
		if strutil.IsNotEmpty(f.environment) {
			env = enum.Environment(enum.Environment_value[f.environment])
		}
		bc.Environment = env
	}
}
