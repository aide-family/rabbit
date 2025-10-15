// Package get is the get command for the Rabbit service
package get

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

func NewCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a message from the queue",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.MessageCommands,
		},
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("get called")
			globalFlags := flags.GlobalFlags
			fmt.Printf("globalFlags: %+v\n", globalFlags)
		},
	}
	flags.addFlags(getCmd)
	return getCmd
}
