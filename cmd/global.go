package cmd

import (
	"os"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
)

var hostname, _ = os.Hostname()

type GlobalFlags struct {
	Helper      *klog.Helper `json:"-" yaml:"-"`
	Name        string       `json:"name" yaml:"name"`
	Author      string       `json:"author" yaml:"author"`
	Email       string       `json:"email" yaml:"email"`
	Repo        string       `json:"repo" yaml:"repo"`
	Description string       `json:"description" yaml:"description"`
	Version     string       `json:"version" yaml:"version"`
	Built       string       `json:"built" yaml:"built"`
	Hostname    string       `json:"-" yaml:"-"`

	Namespace string `json:"-" yaml:"-"`

	RabbitConfigPath string `json:"-" yaml:"-"`
}

func (g *GlobalFlags) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&g.Namespace, "namespace", "n", "", "The namespace of the service")
	cmd.PersistentFlags().StringVar(&g.RabbitConfigPath, "rabbit-config", "~/.rabbit", "The config file of the rabbit")
}

type GlobalOption func(*GlobalFlags)

var globalFlags GlobalFlags = GlobalFlags{
	Name:        "moon.rabbit",
	Author:      "",
	Email:       "",
	Repo:        "https://github.com/aide-family/rabbit",
	Description: "",
	Hostname:    hostname,
	Helper:      klog.NewHelper(klog.DefaultLogger),
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

func WithGlobalFlagsEmail(email string) GlobalOption {
	return func(g *GlobalFlags) {
		g.Email = email
	}
}

func WithGlobalFlagsAuthor(author string) GlobalOption {
	return func(g *GlobalFlags) {
		g.Author = author
	}
}

func WithGlobalFlagsDescription(description string) GlobalOption {
	return func(g *GlobalFlags) {
		g.Description = description
	}
}

func WithGlobalFlagsREPO(repo string) GlobalOption {
	return func(g *GlobalFlags) {
		g.Repo = repo
	}
}

func WithGlobalFlagsHelper(helper *klog.Helper) GlobalOption {
	return func(g *GlobalFlags) {
		g.Helper = helper
	}
}
