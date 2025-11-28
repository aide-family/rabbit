/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	_ "embed"

	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/apply"
	"github.com/aide-family/rabbit/cmd/config"
	"github.com/aide-family/rabbit/cmd/delete"
	"github.com/aide-family/rabbit/cmd/get"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/cmd/send"
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
	cmd.SetGlobalFlags(
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
		run.NewCmd(defaultServerConfig),
		send.NewCmd(),
		version.NewCmd(),
	}
	cmd.Execute(rootCmd, children...)
}
