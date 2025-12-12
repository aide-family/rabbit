// Package run is the run command for the Rabbit service
package run

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/internal/conf"
)

const cmdRunLong = `Run the Rabbit services`

func NewCmd(defaultServerConfigBytes []byte) *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Rabbit services",
		Long:  cmdRunLong,
	}
	var bc conf.Bootstrap
	c := config.New(config.WithSource(
		env.NewSource(),
		conf.NewBytesSource(defaultServerConfigBytes),
	), config.WithPrintLoadedDebugLog(false))
	if err := c.Load(); err != nil {
		klog.Errorw("msg", "load config failed", "error", err)
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		klog.Errorw("msg", "scan config failed", "error", err)
		panic(err)
	}
	runFlags.addFlags(runCmd, &bc)

	return runCmd
}
