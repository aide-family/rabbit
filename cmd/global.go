package cmd

import (
	"os"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
)

var hostname, _ = os.Hostname()

type GlobalFlags struct {
	Helper   *klog.Helper `json:"-" yaml:"-"`
	Name     string       `json:"name" yaml:"name"`
	Author   string       `json:"author" yaml:"author"`
	Github   string       `json:"github" yaml:"github"`
	Version  string       `json:"version" yaml:"version"`
	Built    string       `json:"built" yaml:"built"`
	Hostname string       `json:"-" yaml:"-"`

	Namespace string `json:"-" yaml:"-"`

	RabbitConfigPath string `json:"-" yaml:"-"`
}

func (g *GlobalFlags) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&g.Namespace, "namespace", "n", "", "The namespace of the service")
	cmd.PersistentFlags().StringVar(&g.RabbitConfigPath, "rabbit-config", "~/.rabbit", "The config file of the rabbit")
}

type GlobalOption func(*GlobalFlags)

var globalFlags GlobalFlags = GlobalFlags{
	Author:   "Aide Family",
	Github:   "https://github.com/aide-family/rabbit",
	Hostname: hostname,
	Helper:   klog.NewHelper(klog.DefaultLogger),
}

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

func WithGlobalFlagsHelper(helper *klog.Helper) GlobalOption {
	return func(g *GlobalFlags) {
		g.Helper = helper
	}
}
