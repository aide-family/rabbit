package send

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/rabbit/cmd"
)

type SendFlags struct {
	*cmd.GlobalFlags
}

var sendFlags SendFlags

func (f *SendFlags) addFlags(c *cobra.Command) {
	f.GlobalFlags = cmd.GetGlobalFlags()
}

func GetSendFlags() SendFlags {
	if pointer.IsNil(sendFlags.GlobalFlags) {
		sendFlags.GlobalFlags = cmd.GetGlobalFlags()
	}

	return sendFlags
}
