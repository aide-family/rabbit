package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var hostname, _ = os.Hostname()

type GlobalFlags struct {
	Name     string `json:"name" yaml:"name"`
	Author   string `json:"author" yaml:"author"`
	Github   string `json:"github" yaml:"github"`
	Version  string `json:"version" yaml:"version"`
	Built    string `json:"built" yaml:"built"`
	Hostname string `json:"-" yaml:"-"`

	Namespace string `json:"-" yaml:"-"`

	rabbitConfigPath string `json:"-" yaml:"-"`
}

func (g *GlobalFlags) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&g.Name, "name", "Rabbit", "The name of the service")
	cmd.PersistentFlags().StringVar(&g.Author, "author", "Aide Family", "The author of the service")
	cmd.PersistentFlags().StringVar(&g.Github, "github", "https://github.com/aide-family/rabbit", "The github of the service")
	cmd.PersistentFlags().StringVar(&g.Hostname, "hostname", hostname, "The hostname of the service")
	cmd.PersistentFlags().StringVarP(&g.Namespace, "namespace", "n", "", "The namespace of the service")
	cmd.PersistentFlags().StringVar(&g.rabbitConfigPath, "rabbit-config", "~/.rabbit/config", "The config file of the rabbit")
}

type GlobalOption func(*GlobalFlags)

var globalFlags GlobalFlags

func GetGlobalFlags() GlobalFlags {
	return globalFlags
}

func SetGlobalFlags(opts ...GlobalOption) {
	for _, opt := range opts {
		opt(&globalFlags)
	}
}

func WithGlobalFlagsVersion(version string) GlobalOption {
	return func(g *GlobalFlags) {
		g.Version = version
	}
}

func WithGlobalFlagsBuildTime(buildTime string) GlobalOption {
	return func(g *GlobalFlags) {
		g.Built = buildTime
	}
}
