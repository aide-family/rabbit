package data

import (
	sync "sync"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
)

const (
	KeyNamespaces = "namespaces"
	KeyWebhooks   = "webhooks"
	KeyEmails     = "emails"
	KeyTemplates  = "templates"
)

var (
	keys           = []string{KeyNamespaces, KeyWebhooks, KeyEmails, KeyTemplates}
	fileConfigOnce sync.Once
)

// GetFileConfig returns the file-based configuration
func (d *Data) GetFileConfig() *conf.Config {
	return &d.fileConfig
}

func (d *Data) RegisterReloadFunc(key string, fn func()) {
	d.reloadFuncs.Set(key, fn)
}

// LoadFileConfig loads configuration from configPaths directories using kratos config system
func (d *Data) LoadFileConfig(bc *conf.Bootstrap, helper *klog.Helper) error {
	reloadFunc := func(c config.Config, key string) {
		if err := c.Scan(&d.fileConfig); err != nil {
			helper.Errorw("msg", "scan config failed", "error", err)
			return
		}
		if reloadFunc, ok := d.reloadFuncs.Get(key); ok {
			reloadFunc()
		}
	}
	var err error
	fileConfigOnce.Do(func() {
		if bc.GetUseDatabase() == "true" {
			helper.Debugw("msg", "database mode is enabled, skipping file config loading")
			return
		}
		configPaths := bc.GetConfigPaths()
		if len(configPaths) == 0 {
			helper.Debugw("msg", "no configPaths specified, skipping file config loading")
			return
		}

		// Collect all file paths
		fileSources := make([]config.Source, 0, len(configPaths))
		fileSources = append(fileSources, env.NewSource())
		for _, configPath := range strutil.SplitSkipEmpty(configPaths, ",") {
			fileSources = append(fileSources, file.NewSource(load.ExpandHomeDir(configPath)))
		}

		if len(fileSources) == 0 {
			helper.Debugw("msg", "no config files found in configPaths")
			return
		}

		c := config.New(config.WithSource(fileSources...))
		if err = c.Load(); err != nil {
			helper.Errorw("msg", "load config failed", "error", err)
			return
		}

		if err = c.Scan(&d.fileConfig); err != nil {
			helper.Errorw("msg", "scan config failed", "error", err)
			return
		}
		for _, key := range keys {
			c.Watch(key, func(key string, value config.Value) {
				helper.Debugw("msg", "file config changed", "key", key, "value", value)
				reloadFunc(c, key)
			})
		}
	})

	return err
}
