/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	"github.com/spf13/cobra"

	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/apply"
	"github.com/aide-family/rabbit/cmd/delete"
	"github.com/aide-family/rabbit/cmd/get"
	"github.com/aide-family/rabbit/cmd/gorm"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/cmd/send"
	"github.com/aide-family/rabbit/cmd/version"
)

var (
	Version          = "latest"
	BuildTime string = "2025-10-14T10:00:00Z"
)

func main() {
	logger, err := log.NewLogger(stdio.LoggerDriver())
	if err != nil {
		panic(err)
	}
	logger = klog.With(logger,
		"ts", klog.DefaultTimestamp,
	)
	helper := klog.NewHelper(logger)
	cmd.SetGlobalFlags(
		cmd.WithGlobalFlagsHelper(helper),
		cmd.WithGlobalFlagsVersion(Version),
		cmd.WithGlobalFlagsBuildTime(BuildTime),
	)
	rootCmd := cmd.NewCmd()

	children := []*cobra.Command{
		apply.NewCmd(),
		delete.NewCmd(),
		get.NewCmd(),
		run.NewCmd(),
		send.NewCmd(),
		version.NewCmd(),
		gorm.NewCmd(),
	}
	cmd.Execute(rootCmd, children...)
}
