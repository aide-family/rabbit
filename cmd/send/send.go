// Package send is the send command for the Rabbit service
package send

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/send/email"
	"github.com/aide-family/rabbit/cmd/send/feishu"
	"github.com/aide-family/rabbit/cmd/send/sms"
)

func NewCmd() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send a message to the queue",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	commands := []*cobra.Command{
		sms.NewCmd(),
		feishu.NewCmd(),
		email.NewCmd(),
	}
	sendCmd.AddCommand(commands...)

	flags.addFlags(sendCmd)

	return sendCmd
}
