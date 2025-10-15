// Package email is the Email command for the Rabbit service
package email

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

func NewCmd() *cobra.Command {
	emailCmd := &cobra.Command{
		Use:   "email",
		Short: "Send an email message",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
	flags.addFlags(emailCmd)
	return emailCmd
}
