/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

// Package cmd is the root command for the Rabbit service
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Command groups for organized help display
const (
	BasicCommands   = "Basic Commands"
	MessageCommands = "Message Commands"
	ServiceCommands = "Service Commands"
)

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "rabbit",
		Short: "Moon messaging platform - Rabbit service",
		Long: `Rabbit（玉兔）是 Moon 平台的消息服务工具，提供统一的消息发送和管理能力。

Rabbit 是一个基于 Kratos 框架构建的分布式消息服务平台，支持多种消息通道
（邮件、Webhook、短信、飞书等）的统一管理和发送。通过命名空间（Namespace）
实现多租户隔离，支持配置文件和数据库两种存储模式，满足不同场景的部署需求。

核心能力：
  • 多通道消息发送：支持邮件、Webhook、短信、飞书等多种消息通道的统一管理
  • 模板化发送：支持消息模板配置，实现消息内容的动态渲染和复用
  • 异步消息处理：基于消息队列实现异步发送，提升系统吞吐量和可靠性
  • 配置管理：支持邮件服务器、Webhook 端点等通道配置的集中管理
  • 多租户隔离：通过命名空间实现不同业务或租户的配置和数据隔离
  • 命令行工具：提供丰富的 CLI 命令，支持服务管理、消息发送、配置生成等

命令分类：
  • Basic Commands（基础命令）：config、version 等基础操作
  • Message Commands（消息命令）：send、apply、get、delete 等消息相关操作
  • Service Commands（服务命令）：run 等服务管理操作

使用场景：
  • 企业级通知系统：统一管理各类业务通知（订单、告警、系统消息等）
  • 微服务消息中心：为微服务架构提供统一的消息发送能力
  • 多渠道推送平台：集成多种消息通道，实现消息的统一发送和管理
  • 开发测试工具：通过命令行快速测试消息通道配置和发送功能

使用 "rabbit [command] --help" 查看具体命令的详细说明。`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	globalFlags.addFlags(rootCmd)
	// Set custom help template to display commands in groups
	rootCmd.SetHelpTemplate(customHelpTemplate)
	rootCmd.SetUsageTemplate(customUsageTemplate)

	// Register custom template function
	cobra.AddTemplateFunc("customCommands", func(cmd *cobra.Command) string {
		return Commands(cmd)
	})

	return rootCmd
}

// customHelpTemplate is the custom help template that groups commands
var customHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// customUsageTemplate provides custom usage formatting with command groups
var customUsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
{{. | customCommands}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// Commands returns the grouped commands for help display
func Commands(cmd *cobra.Command) string {
	groups := make(map[string][]*cobra.Command)

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		group := getCommandGroup(c)
		groups[group] = append(groups[group], c)
	}

	// Define group order
	groupOrder := []string{BasicCommands, MessageCommands, ServiceCommands}

	var result strings.Builder
	for _, groupName := range groupOrder {
		if commands, exists := groups[groupName]; exists {
			sort.Slice(commands, func(i, j int) bool {
				return commands[i].Name() < commands[j].Name()
			})

			result.WriteString(fmt.Sprintf("\n%s:\n", groupName))
			for _, c := range commands {
				result.WriteString(fmt.Sprintf("  %-15s %s\n", c.Name(), c.Short))
			}
		}
	}

	// Add any remaining commands that don't have a group
	for groupName, commands := range groups {
		found := false
		for _, orderedGroup := range groupOrder {
			if groupName == orderedGroup {
				found = true
				break
			}
		}
		if !found {
			sort.Slice(commands, func(i, j int) bool {
				return commands[i].Name() < commands[j].Name()
			})
			result.WriteString(fmt.Sprintf("\n%s:\n", groupName))
			for _, c := range commands {
				result.WriteString(fmt.Sprintf("  %-15s %s\n", c.Name(), c.Short))
			}
		}
	}

	return result.String()
}

// getCommandGroup determines which group a command belongs to
func getCommandGroup(cmd *cobra.Command) string {
	// Check if command has an annotation for its group
	if group, exists := cmd.Annotations["group"]; exists {
		return group
	}

	return BasicCommands
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cmd *cobra.Command, children ...*cobra.Command) {
	cmd.AddCommand(children...)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
