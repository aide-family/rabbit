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

const txtTemplate = `Name:	{{.Name}}
Author:	{{.Author}}
Email:	{{.Email}}
Version:{{.Version}}
Repo:	{{.Repo}}
Built:	{{.Built}}
Description:	{{.Description}}
`

const cmdLong = `Display version information and build details for the Rabbit service.

The version command displays detailed information about the current Rabbit service,
including version number, build time, author information, and other metadata to help
understand the service version and manage versioning.

Key Features:
  • Version information: Display basic information such as version number and build time
  • Project metadata: Display project name, author, repository URL, and other metadata
  • Multiple output formats: Support for text, JSON, YAML, and other output formats
  • Detailed description: Display project functionality description and usage instructions

Output Formats:
  • Default format: Display version information in human-readable text format
  • JSON format: Use --format json to output in JSON format
  • YAML format: Use --format yaml to output in YAML format

Use Cases:
  • Version verification: Confirm the version of the currently running Rabbit service
  • Environment validation: Verify that the service version in the deployment environment is correct
  • Troubleshooting: Include version information in issue reports to facilitate problem diagnosis

Version information is crucial for troubleshooting and version management. It is recommended
to record version information during deployment and when reporting issues.`

func NewCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information and build details for the Rabbit service",
		Long:  cmdLong,
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
