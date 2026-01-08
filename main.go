/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	_ "embed"
	"os"

	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/cmd/run/all"
	"github.com/aide-family/rabbit/cmd/run/grpc"
	"github.com/aide-family/rabbit/cmd/run/http"
	"github.com/aide-family/rabbit/cmd/run/job"
	"github.com/aide-family/rabbit/cmd/version"
	"github.com/aide-family/rabbit/pkg/merr"
)

var (
	Name        = "moon"
	Version     = "latest"
	BuildTime   = "now"
	Author      = "Aide Family"
	Email       = "aidecloud@163.com"
	Repo        = "https://github.com/aide-family/rabbit"
	hostname, _ = os.Hostname()
)

//go:embed description.txt
var Description string

//go:embed config/server.yaml
var defaultServerConfig []byte

func main() {
	cmd.SetGlobalFlags(
		cmd.WithGlobalFlagsName(Name),
		cmd.WithGlobalFlagsHostname(hostname),
		cmd.WithGlobalFlagsVersion(Version),
		cmd.WithGlobalFlagsBuildTime(BuildTime),
		cmd.WithGlobalFlagsAuthor(Author),
		cmd.WithGlobalFlagsEmail(Email),
		cmd.WithGlobalFlagsREPO(Repo),
		cmd.WithGlobalFlagsDescription(Description),
	)

	runCmd := run.NewCmd(defaultServerConfig)
	runCmd.AddCommand(grpc.NewCmd(), http.NewCmd(), job.NewCmd(), all.NewCmd())

	children := []*cobra.Command{
		runCmd,
		version.NewCmd(),
	}
	cmd.Execute(cmd.NewCmd(), children...)
}

func init() {
	logger, err := log.NewLogger(stdio.LoggerDriver())
	if err != nil {
		panic(merr.ErrorInternal("new logger failed with error: %v", err).WithCause(err))
	}
	logger = klog.With(logger,
		"ts", klog.DefaultTimestamp,
	)
	filterLogger := klog.NewFilter(logger, klog.FilterLevel(klog.LevelInfo))
	helper := klog.NewHelper(filterLogger)
	klog.SetLogger(helper.Logger())
}
