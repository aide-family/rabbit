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

const cmdConfigLong = `Generate default configuration file for Rabbit service to enable quick initialization and deployment.

The config command generates a server.yaml configuration file at the specified path, containing
all default configuration items required for service operation, such as server ports, database
connections, log levels, etc.

Key Features:
  • Default configuration generation: Automatically generates configuration files with all required settings
  • Safe backup: Automatically backs up existing files with timestamp suffix if target file already exists
  • Force overwrite: Supports --force flag to forcefully overwrite existing configuration files
  • Customizable path: Supports specifying the generation path and filename for configuration files

Use Cases:
  • Quick initialization: Rapidly generate configuration files for new deployment environments
  • Configuration template: Use as a configuration template reference to understand all configurable options
  • Configuration migration: Serve as a starting point when migrating configurations between different environments

The generated configuration file includes settings for all modules such as server, database, logging,
and message channels, which can be modified and adjusted according to the actual environment.`

func NewCmd(configContent []byte) *cobra.Command {
	defaultConfigContent = configContent
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Generate default server.yaml configuration file",
		Long:  cmdConfigLong,
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
