package server

import (
	"context"
	"net/url"
	"strings"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/service"
)

var _ transport.Server = (*Job)(nil)

type JobServer interface {
	transport.Server
	transport.Endpointer
}

func NewJob(bc *conf.Bootstrap, namespaceService *service.NamespaceService, helper *klog.Helper) *Job {
	return newJob(bc.GetServer().GetJob(), bc.GetJwt(), namespaceService, helper)
}

func newJob(jobConf conf.ServerConfig, jwtConf conf.JWTConfig, namespaceService *service.NamespaceService, helper *klog.Helper) *Job {
	var server JobServer
	serverConf := &conf.Server_ServerConfig{
		Address: jobConf.GetAddress(),
		Timeout: jobConf.GetTimeout(),
	}
	if strings.EqualFold(jobConf.GetNetwork(), string(transport.KindHTTP)) {
		server = newHTTPServer(serverConf, jwtConf, namespaceService, helper)
	} else {
		server = newGRPCServer(serverConf, jwtConf, namespaceService, helper)
	}
	return &Job{
		helper: klog.NewHelper(klog.With(helper.Logger(), "server", "job")),
		server: server,
	}
}

type Job struct {
	jobService *service.JobService
	server     JobServer
	helper     *klog.Helper
}

func (e *Job) RegisterHandler(jobService *service.JobService) {
	e.jobService = jobService
}

// Start implements transport.Server.
func (e *Job) Start(ctx context.Context) error {
	if err := e.jobService.Start(ctx); err != nil {
		e.helper.Errorw("msg", "start job failed", "error", err)
		return err
	}
	endpoint, err := e.server.Endpoint()
	if err != nil {
		e.helper.Errorw("msg", "get job endpoint failed", "error", err)
		return err
	}
	e.helper.Infow("msg", "[Job] started", "address", endpoint)
	if err := e.server.Start(ctx); err != nil {
		e.helper.Errorw("msg", "start server failed", "error", err)
		return err
	}
	return nil
}

// Stop implements transport.Server.
func (e *Job) Stop(ctx context.Context) error {
	if err := e.server.Stop(ctx); err != nil {
		e.helper.Errorw("msg", "stop server failed", "error", err)
		return err
	}
	if err := e.jobService.Stop(ctx); err != nil {
		e.helper.Errorw("msg", "stop job failed", "error", err)
		return err
	}
	e.helper.Infow("msg", "[Job] stopped")
	return nil
}

// Endpoint implements transport.Server.
func (e *Job) Endpoint() (*url.URL, error) {
	return e.server.Endpoint()
}
