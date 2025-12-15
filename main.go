/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/apply"
	"github.com/aide-family/rabbit/cmd/config"
	"github.com/aide-family/rabbit/cmd/delete"
	"github.com/aide-family/rabbit/cmd/get"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/cmd/run/all"
	"github.com/aide-family/rabbit/cmd/run/grpc"
	"github.com/aide-family/rabbit/cmd/run/http"
	"github.com/aide-family/rabbit/cmd/run/job"
	"github.com/aide-family/rabbit/cmd/send"
	"github.com/aide-family/rabbit/cmd/send/email"
	"github.com/aide-family/rabbit/cmd/send/feishu"
	"github.com/aide-family/rabbit/cmd/send/sms"
	"github.com/aide-family/rabbit/cmd/version"
)

var (
	Name        = "moon.rabbit"
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

	sendCmd := send.NewCmd(sms.NewCmd(), feishu.NewCmd(), email.NewCmd())
	runCmd := run.NewCmd(defaultServerConfig)
	runCmd.AddCommand(grpc.NewCmd(), http.NewCmd(), job.NewCmd(), all.NewCmd())

	children := []*cobra.Command{
		apply.NewCmd(),
		config.NewCmd(defaultServerConfig),
		delete.NewCmd(),
		get.NewCmd(),
		sendCmd,
		runCmd,
		version.NewCmd(),
	}
	cmd.Execute(cmd.NewCmd(), children...)
}
