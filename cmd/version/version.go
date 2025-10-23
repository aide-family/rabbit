// Package version is the version command for the Rabbit service
package version

import (
	"fmt"
	"os"
	"text/template"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

var txtTemplate = `Name:	{{.Name}}
Author:	{{.Author}}
Version:{{.Version}}
Github:	{{.Github}}
Built:	{{.Built}}
`

func NewCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of the Rabbit service",
		Long:  `Show the version of the Rabbit service`,
		Annotations: map[string]string{
			"group": cmd.BasicCommands,
		},
		Run: func(c *cobra.Command, args []string) {
			flags.GlobalFlags = cmd.GetGlobalFlags()
			switch flags.format {
			case "json", "yaml":
				bytes, _ := encoding.GetCodec(flags.format).Marshal(flags.GlobalFlags)
				fmt.Println(string(bytes))
			default:
				t := template.Must(template.New("txt").Parse(txtTemplate))
				t.Execute(os.Stdout, flags.GlobalFlags)
			}
		},
	}
	flags.addFlags(versionCmd)
	return versionCmd
}
