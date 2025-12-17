package server

import (
	"context"
	"net/url"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/service"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/merr"
)

var (
	_ transport.Server     = (*JobServer)(nil)
	_ transport.Endpointer = (*JobServer)(nil)
)

func NewJobServer(bc *conf.Bootstrap, namespaceService *service.NamespaceService, helper *klog.Helper) (*JobServer, error) {
	return newJobServer(bc.GetServer().GetJob(), bc.GetJwt(), namespaceService, helper)
}

func newJobServer(jobConf conf.ServerConfig, jwtConf conf.JWTConfig, namespaceService *service.NamespaceService, helper *klog.Helper) (*JobServer, error) {
	protocol := jobConf.GetProtocol()
	job := &JobServer{
		helper:   klog.NewHelper(klog.With(helper.Logger(), "server", "job")),
		protocol: protocol,
	}

	switch protocol {
	case config.ClusterConfig_HTTP:
		job.httpSrv = newHTTPServer(jobConf, jwtConf, namespaceService, helper)
		endpoint, err := job.httpSrv.Endpoint()
		if err != nil {
			helper.Errorw("msg", "get job endpoint failed", "error", err)
			return nil, err
		}
		job.endpoint = endpoint
		job.startFunc = job.httpSrv.Start
		job.stopFunc = job.httpSrv.Stop
	case config.ClusterConfig_GRPC:
		job.grpcSrv = newGRPCServer(jobConf, jwtConf, namespaceService, helper)
		endpoint, err := job.grpcSrv.Endpoint()
		if err != nil {
			helper.Errorw("msg", "get job endpoint failed", "error", err)
			return nil, err
		}
		job.endpoint = endpoint
		job.startFunc = job.grpcSrv.Start
		job.stopFunc = job.grpcSrv.Stop
	default:
		return nil, merr.ErrorInternal("not support protocol: %s", protocol)
	}

	return job, nil
}

type JobServer struct {
	protocol  config.ClusterConfig_Protocol
	httpSrv   *http.Server
	grpcSrv   *grpc.Server
	helper    *klog.Helper
	endpoint  *url.URL
	startFunc func(ctx context.Context) error
	stopFunc  func(ctx context.Context) error
}

// Start implements transport.Server.
func (e *JobServer) Start(ctx context.Context) error {
	e.helper.Infow("msg", "[Job] started", "address", e.endpoint)
	if err := e.startFunc(ctx); err != nil {
		e.helper.Errorw("msg", "start server failed", "error", err)
		return err
	}
	return nil
}

// Stop implements transport.Server.
func (e *JobServer) Stop(ctx context.Context) error {
	if err := e.stopFunc(ctx); err != nil {
		e.helper.Errorw("msg", "stop server failed", "error", err)
		return err
	}
	e.helper.Infow("msg", "[Job] stopped")
	return nil
}

// Endpoint implements transport.Server.
func (e *JobServer) Endpoint() (*url.URL, error) {
	return e.endpoint, nil
}
