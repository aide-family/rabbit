// Package delete is the delete command for the Rabbit service
package delete

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

func NewCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a message from the queue",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
		},
	}
	flags.addFlags(deleteCmd)
	return deleteCmd
}
