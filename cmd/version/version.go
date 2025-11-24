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
Email:	{{.Email}}
Version:{{.Version}}
Repo:	{{.Repo}}
Built:	{{.Built}}
Description:	{{.Description}}
`

func NewCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of the Rabbit service",
		Long: `显示 Rabbit 服务的版本信息和构建详情。

version 命令用于查看当前 Rabbit 服务的版本号、构建时间、作者信息
等详细信息，帮助了解服务版本和进行版本管理。

主要功能：
  • 版本信息：显示服务的版本号、构建时间等基本信息
  • 项目信息：显示项目名称、作者、仓库地址等元数据
  • 多格式输出：支持文本、JSON、YAML 等多种格式输出
  • 详细描述：显示项目的功能描述和使用说明

输出格式：
  • 默认格式：以易读的文本格式显示版本信息
  • JSON 格式：使用 --format json 输出 JSON 格式
  • YAML 格式：使用 --format yaml 输出 YAML 格式

使用场景：
  • 版本检查：确认当前运行的 Rabbit 服务版本
  • 环境验证：验证部署环境中的服务版本是否正确
  • 问题排查：在问题报告中包含版本信息，便于定位问题

版本信息对于问题排查和版本管理非常重要，建议在部署和问题报告时
记录版本信息。`,
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
