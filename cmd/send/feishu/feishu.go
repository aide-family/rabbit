// Package feishu is the Feishu command for the Rabbit service
package feishu

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "feishu",
		Short: "Send a Feishu message",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("feishu called")
}
