package main

import (
	"strings"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/config"
)

type Flags struct {
	*cmd.GlobalFlags
	configPath string
	forceGen   bool
	isBiz      bool

	ormConf *config.ORMConfig
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.GlobalFlags = cmd.GetGlobalFlags()
	c.PersistentFlags().StringVarP(&f.configPath, "config", "c", "./config", "config file")
	c.PersistentFlags().BoolVarP(&f.forceGen, "force-gen", "f", false, "force generate code, if the code already exists, it will be overwitten")
	c.PersistentFlags().BoolVarP(&f.isBiz, "biz", "b", false, "is biz")
}

// containsNamespace 检查 namespace 是否在逗号分隔的 namespaces 字符串中
func containsNamespace(namespaces, target string) bool {
	for _, ns := range strutil.SplitSkipEmpty(namespaces, ",") {
		if strings.TrimSpace(ns) == target {
			return true
		}
	}
	return false
}

func (f *Flags) applyToBootstrap(bc *conf.Bootstrap) {
	if pointer.IsNil(bc) || pointer.IsNil(bc.GetMain()) {
		return
	}
	ormConf := bc.GetMain()
	if f.isBiz {
		for namespaces, bizConfig := range bc.GetBiz() {
			if containsNamespace(namespaces, f.Namespace) {
				ormConf = bizConfig
				break
			}
		}
	}
	f.ormConf = ormConf
}
