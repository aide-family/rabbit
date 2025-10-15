// Package apply is the apply command for the Rabbit service
package apply

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

func NewCmd() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a message to the queue",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply called")
		},
	}
	flags.addFlags(applyCmd)
	return applyCmd
}
