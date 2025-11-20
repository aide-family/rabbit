// Package conf provides configuration for the application.
package conf

import (
	sync "sync"

	"github.com/aide-family/magicbox/load"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
)

var (
	fileConfig     Config
	fileConfigOnce sync.Once
)

// GetFileConfig returns the file-based configuration
func GetFileConfig() *Config {
	return &fileConfig
}

// LoadFileConfig loads configuration from configPaths directories using kratos config system
func LoadFileConfig(bc *Bootstrap, helper *klog.Helper) error {
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
	})

	return err
}
