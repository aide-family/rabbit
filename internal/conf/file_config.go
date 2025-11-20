// Package conf provides configuration for the application.
package conf

import (
	sync "sync"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/safety"
	"github.com/go-kratos/kratos/v2/config"
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
	fileConfig     Config
	fileConfigOnce sync.Once
	reloadFuncs    = safety.NewSyncMap(make(map[string]func()))
)

// GetFileConfig returns the file-based configuration
func GetFileConfig() *Config {
	return &fileConfig
}

func RegisterReloadFunc(key string, fn func()) {
	reloadFuncs.Set(key, fn)
}

// LoadFileConfig loads configuration from configPaths directories using kratos config system
func LoadFileConfig(bc *Bootstrap, helper *klog.Helper) error {
	reloadFunc := func(c config.Config, key string) {
		if err := c.Scan(&fileConfig); err != nil {
			helper.Errorw("msg", "scan config failed", "error", err)
			return
		}
		if reloadFunc, ok := reloadFuncs.Get(key); ok {
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
		var fileSources []config.Source
		for _, configPath := range configPaths {
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

		if err = c.Scan(&fileConfig); err != nil {
			helper.Errorw("msg", "scan config failed", "error", err)
			return
		}
		c.Watch(KeyNamespaces, func(key string, value config.Value) {
			helper.Infow("msg", "namespaces changed", "key", key, "value", value)
			reloadFunc(c, key)
		})
		c.Watch(KeyWebhooks, func(key string, value config.Value) {
			helper.Infow("msg", "webhooks changed", "key", key, "value", value)
			reloadFunc(c, key)
		})
		c.Watch(KeyEmails, func(key string, value config.Value) {
			helper.Infow("msg", "emails changed", "key", key, "value", value)
			reloadFunc(c, key)
		})
		c.Watch(KeyTemplates, func(key string, value config.Value) {
			helper.Infow("msg", "templates changed", "key", key, "value", value)
			reloadFunc(c, key)
		})
	})

	return err
}
