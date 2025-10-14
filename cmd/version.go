/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/spf13/cobra"
)

type ServerInfo struct {
	Name      string `json:"name" yaml:"name"`
	Author    string `json:"author" yaml:"author"`
	Version   string `json:"version" yaml:"version"`
	BuildTime string `json:"build_time" yaml:"build_time"`
	Github    string `json:"github" yaml:"github"`
}

var format = "txt"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the Rabbit service",
	Long:  `Show the version of the Rabbit service`,
	Run: func(cmd *cobra.Command, args []string) {
		serverInfo := ServerInfo{
			Name:      "Rabbit",
			Author:    "Aide Family",
			Version:   Version,
			BuildTime: BuildTime,
			Github:    "https://github.com/aide-family/rabbit",
		}

		switch format {
		case "json", "yaml":
			bytes, _ := encoding.GetCodec(format).Marshal(serverInfo)
			fmt.Println(string(bytes))
		default:
			txt := `Name: %s
Author: %s
Version: %s
BuildTime: %s
Github: %s
`
			fmt.Printf(txt, serverInfo.Name, serverInfo.Author, serverInfo.Version, serverInfo.BuildTime, serverInfo.Github)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")
	versionCmd.PersistentFlags().StringVarP(&format, "format", "f", "txt", "The format of the version output, supported: txt, json, yaml")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
