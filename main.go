/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	"github.com/aide-family/rabbit/cmd"
)

var (
	Version   = "latest"
	BuildTime string
)

func main() {
	cmd.Execute(Version, BuildTime)
}
