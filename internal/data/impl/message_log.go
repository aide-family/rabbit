package impl

import (
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/dbimpl"
	"github.com/aide-family/rabbit/internal/data/impl/fileimpl"
)

func NewMessageLogRepository(d *data.Data, helper *klog.Helper) repository.MessageLog {
	newRepo := fileimpl.NewMessageLogRepository
	if d.UseDatabase() {
		newRepo = dbimpl.NewMessageLogRepository
	}
	return newRepo(d, helper)
}
