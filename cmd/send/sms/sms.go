// Package sms is the SMS command for the Rabbit service
package sms

import (
	"fmt"

	"github.com/aide-family/rabbit/cmd"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sms",
		Short: "Send a SMS message",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("sms called")
}
