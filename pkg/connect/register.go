package connect

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/rabbit/pkg/config"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var globalRegistry = NewRegistry()

type ORMFactory func(c *config.ORMConfig, logger *klog.Helper) (*gorm.DB, error)
type ReportFactory func(c *config.ReportConfig, logger *klog.Helper) (Report, func() error, error)

type Registry struct {
	ormConfigs    *safety.SyncMap[config.ORMConfig_Dialector, ORMFactory]
	reportConfigs *safety.SyncMap[config.ReportConfig_ReportType, ReportFactory]
}

func (r *Registry) RegisterORMFactory(dialector config.ORMConfig_Dialector, factory ORMFactory) {
	r.ormConfigs.Set(dialector, factory)
}

func (r *Registry) GetORMFactory(dialector config.ORMConfig_Dialector) (ORMFactory, bool) {
	return r.ormConfigs.Get(dialector)
}

func (r *Registry) RegisterReportFactory(registryType config.ReportConfig_ReportType, factory ReportFactory) {
	r.reportConfigs.Set(registryType, factory)
}

func (r *Registry) GetReportFactory(registryType config.ReportConfig_ReportType) (ReportFactory, bool) {
	return r.reportConfigs.Get(registryType)
}

func NewRegistry() *Registry {
	return &Registry{
		ormConfigs:    safety.NewSyncMap(make(map[config.ORMConfig_Dialector]ORMFactory)),
		reportConfigs: safety.NewSyncMap(make(map[config.ReportConfig_ReportType]ReportFactory)),
	}
}
