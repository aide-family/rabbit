package main

import (
	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	"github.com/aide-family/rabbit/pkg/merr"
	klog "github.com/go-kratos/kratos/v2/log"
)

func init() {
	logger, err := log.NewLogger(stdio.LoggerDriver())
	if err != nil {
		panic(merr.ErrorInternal("new logger failed with error: %v", err).WithCause(err))
	}
	logger = klog.With(logger,
		"ts", klog.DefaultTimestamp,
	)
	filterLogger := klog.NewFilter(logger, klog.FilterLevel(klog.LevelInfo))
	helper := klog.NewHelper(filterLogger)
	klog.SetLogger(helper.Logger())
}
