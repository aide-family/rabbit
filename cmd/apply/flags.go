package apply

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

type Flags struct {
	cmd.GlobalFlags
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.GlobalFlags = cmd.GetGlobalFlags()
}
