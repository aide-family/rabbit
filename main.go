/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	_ "embed"

	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/apply"
	"github.com/aide-family/rabbit/cmd/config"
	"github.com/aide-family/rabbit/cmd/delete"
	"github.com/aide-family/rabbit/cmd/get"
	"github.com/aide-family/rabbit/cmd/job"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/cmd/send"
	"github.com/aide-family/rabbit/cmd/server"
	"github.com/aide-family/rabbit/cmd/version"
)

var (
	Version   = "latest"
	BuildTime = "now"
	Author    = "Aide Family"
	Email     = ""
	Repo      = "https://github.com/aide-family/rabbit"
)

//go:embed description.txt
var Description string

//go:embed config/server.yaml
var defaultServerConfig []byte

func main() {
	logger, err := log.NewLogger(stdio.LoggerDriver())
	if err != nil {
		panic(err)
	}
	logger = klog.With(logger,
		"ts", klog.DefaultTimestamp,
	)
	filterLogger := klog.NewFilter(logger, klog.FilterLevel(klog.LevelInfo))
	helper := klog.NewHelper(filterLogger)
	klog.SetLogger(helper.Logger())

	cmd.SetGlobalFlags(
		cmd.WithGlobalFlagsHelper(helper),
		cmd.WithGlobalFlagsVersion(Version),
		cmd.WithGlobalFlagsBuildTime(BuildTime),
		cmd.WithGlobalFlagsAuthor(Author),
		cmd.WithGlobalFlagsEmail(Email),
		cmd.WithGlobalFlagsREPO(Repo),
		cmd.WithGlobalFlagsDescription(Description),
	)
	rootCmd := cmd.NewCmd()

	children := []*cobra.Command{
		apply.NewCmd(),
		config.NewCmd(defaultServerConfig),
		delete.NewCmd(),
		get.NewCmd(),
		job.NewCmd(),
		run.NewCmd(),
		send.NewCmd(),
		server.NewCmd(),
		version.NewCmd(),
	}
	cmd.Execute(rootCmd, children...)
}
