/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

// Package cmd is the root command for the Rabbit service
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   string
	BuildTime string
	id, _     = os.Hostname()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rabbit",
	Short: "Moon messaging platform - Rabbit service",
	Long: `Rabbit is a messaging service for the Moon platform.

This service provides message handling capabilities for the Moon ecosystem,
including message routing, delivery, and management features.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, buildTime string) {
	Version = version
	BuildTime = buildTime
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rabbit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
