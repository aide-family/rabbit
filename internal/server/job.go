package server

import (
	"context"
	"net/url"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/service"
)

var _ Server = (*JobServer)(nil)

type Server interface {
	transport.Server
	transport.Endpointer
}

func NewJobServer(bc *conf.Bootstrap, namespaceService *service.NamespaceService, helper *klog.Helper) *JobServer {
	return newJobServer(bc.GetServer().GetJob(), bc.GetJwt(), namespaceService, helper)
}

func newJobServer(jobConf conf.ServerConfig, jwtConf conf.JWTConfig, namespaceService *service.NamespaceService, helper *klog.Helper) *JobServer {
	return &JobServer{
		helper: klog.NewHelper(klog.With(helper.Logger(), "server", "job")),
		server: newGRPCServer(jobConf, jwtConf, namespaceService, helper),
	}
}

type JobServer struct {
	server *grpc.Server
	helper *klog.Helper
}

// Start implements transport.Server.
func (e *JobServer) Start(ctx context.Context) error {
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
func (e *JobServer) Stop(ctx context.Context) error {
	if err := e.server.Stop(ctx); err != nil {
		e.helper.Errorw("msg", "stop server failed", "error", err)
		return err
	}
	e.helper.Infow("msg", "[Job] stopped")
	return nil
}

// Endpoint implements transport.Server.
func (e *JobServer) Endpoint() (*url.URL, error) {
	return e.server.Endpoint()
}
