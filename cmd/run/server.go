// Package run is the run command for the Rabbit service
package run

import (
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
	if err := conf.Load(&bc, env.NewSource(), conf.NewBytesSource(defaultServerConfigBytes)); err != nil {
		klog.Errorw("msg", "load config failed", "error", err)
		panic(err)
	}
	runFlags.addFlags(runCmd, &bc)

	return runCmd
}
