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

var _ transport.Server = (*EventBus)(nil)

type EventBusServer interface {
	transport.Server
	transport.Endpointer
}

func NewEventBus(bc *conf.Bootstrap, namespaceService *service.NamespaceService, helper *klog.Helper) *EventBus {
	return newEventBus(bc.GetServer().GetEventBus(), bc.GetJwt(), namespaceService, helper)
}

func newEventBus(eventBusConf conf.ServerConfig, jwtConf conf.JWTConfig, namespaceService *service.NamespaceService, helper *klog.Helper) *EventBus {
	var server EventBusServer
	serverConf := &conf.Server_ServerConfig{
		Address: eventBusConf.GetAddress(),
		Timeout: eventBusConf.GetTimeout(),
	}
	if strings.EqualFold(eventBusConf.GetNetwork(), string(transport.KindHTTP)) {
		server = newHTTPServer(serverConf, jwtConf, namespaceService, helper)
	} else {
		server = newGRPCServer(serverConf, jwtConf, namespaceService, helper)
	}
	return &EventBus{
		helper: klog.NewHelper(klog.With(helper.Logger(), "server", "event_bus")),
		server: server,
	}
}

type EventBus struct {
	eventBusService *service.EventBusService
	server          EventBusServer
	helper          *klog.Helper
}

func (e *EventBus) RegisterHandler(eventBusService *service.EventBusService) {
	e.eventBusService = eventBusService
}

// Start implements transport.Server.
func (e *EventBus) Start(ctx context.Context) error {
	if err := e.eventBusService.Start(ctx); err != nil {
		e.helper.Errorw("msg", "start event bus failed", "error", err)
		return err
	}
	endpoint, err := e.server.Endpoint()
	if err != nil {
		e.helper.Errorw("msg", "get event bus endpoint failed", "error", err)
		return err
	}
	e.helper.Infow("msg", "[EventBus] started", "address", endpoint)
	if err := e.server.Start(ctx); err != nil {
		e.helper.Errorw("msg", "start server failed", "error", err)
		return err
	}
	return nil
}

// Stop implements transport.Server.
func (e *EventBus) Stop(ctx context.Context) error {
	if err := e.server.Stop(ctx); err != nil {
		e.helper.Errorw("msg", "stop server failed", "error", err)
		return err
	}
	if err := e.eventBusService.Stop(ctx); err != nil {
		e.helper.Errorw("msg", "stop event bus failed", "error", err)
		return err
	}
	e.helper.Infow("msg", "[EventBus] stopped")
	return nil
}

// Endpoint implements transport.Server.
func (e *EventBus) Endpoint() (*url.URL, error) {
	endpoint, err := e.server.Endpoint()
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Scheme: "moon",
		Host:   endpoint.Host,
	}, nil
}
