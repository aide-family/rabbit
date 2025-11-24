// Package config is the config command for the Rabbit service
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
)

var defaultConfigContent []byte

func NewCmd(configContent []byte) *cobra.Command {
	defaultConfigContent = configContent
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Generate default server.yaml configuration file",
		Long: `生成 Rabbit 服务的默认配置文件，用于快速初始化和部署。

config 命令会在指定路径生成 server.yaml 配置文件，包含服务运行所需
的默认配置项，如服务器端口、数据库连接、日志级别等。

主要功能：
  • 默认配置生成：自动生成包含所有必需配置项的配置文件
  • 安全备份：如果目标文件已存在，自动备份为带时间戳的文件
  • 强制覆盖：支持 --force 参数强制覆盖已存在的配置文件
  • 路径自定义：支持指定配置文件的生成路径和文件名

使用场景：
  • 快速初始化：新部署环境时快速生成配置文件
  • 配置模板：作为配置模板参考，了解所有可配置项
  • 配置迁移：在不同环境间迁移配置时作为起点

生成的配置文件包含服务器、数据库、日志、消息通道等各个模块的配置
项，可根据实际环境进行修改和调整。`,
		Annotations: map[string]string{
			"group": cmd.BasicCommands,
		},
		Run: runConfig,
	}
	flags.addFlags(configCmd)
	return configCmd
}

func runConfig(_ *cobra.Command, _ []string) {
	flags.GlobalFlags = cmd.GetGlobalFlags()

	// 确保路径存在
	if err := os.MkdirAll(flags.path, 0o755); err != nil {
		flags.Helper.Errorw("msg", "failed to create directory", "path", flags.path, "error", err)
		return
	}

	// 构建完整文件路径
	targetPath := filepath.Join(flags.path, flags.name)

	// 检查文件是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		if flags.force {
			// force 模式：直接覆盖
			flags.Helper.Infow("msg", "overwriting existing file", "path", targetPath)
		} else {
			// 非 force 模式：重命名现有文件
			ext := filepath.Ext(flags.name)
			nameWithoutExt := flags.name[:len(flags.name)-len(ext)]
			// 如果 nameWithoutExt 为空（文件名只有扩展名），使用 "config" 作为默认名称
			if nameWithoutExt == "" {
				nameWithoutExt = "config"
			}
			timestamp := time.Now().Format("20060102150405")
			backupName := fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
			backupPath := filepath.Join(flags.path, backupName)

			if err := os.Rename(targetPath, backupPath); err != nil {
				flags.Helper.Errorw("msg", "failed to rename existing file", "old", targetPath, "new", backupPath, "error", err)
				return
			}
			flags.Helper.Infow("msg", "existing file renamed", "old", targetPath, "new", backupPath)
		}
	}

	// 写入新文件
	if err := os.WriteFile(targetPath, defaultConfigContent, 0o644); err != nil {
		flags.Helper.Errorw("msg", "failed to write config file", "path", targetPath, "error", err)
		return
	}

	flags.Helper.Infow("msg", "config file generated successfully", "path", targetPath)
}
