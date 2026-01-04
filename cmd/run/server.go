// Package run is the run command for the Rabbit service
package run

import (
	"strings"
	"sync"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/env"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/cobra"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/server"
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

func NewEndpoint(serviceName string, wireApp WireAppFunc) *endpoint {
	return &endpoint{
		serviceName: strings.Join([]string{runFlags.Name, runFlags.Server.Name, serviceName}, "."),
		wireAppFunc: wireApp,
	}
}

func NewEngine(endpoints ...*endpoint) *Engine {
	return &Engine{
		endpoints:   endpoints,
		beforeFuncs: []func(){hello.Hello},
		afterFuncs:  []func(){},
	}
}

type WireAppFunc func(serviceName string, bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, func(), error)

type endpoint struct {
	serviceName string
	wireAppFunc WireAppFunc
	app         *kratos.App
	cleanup     func()
	helper      *klog.Helper
	err         error
}

type Engine struct {
	endpoints   []*endpoint
	beforeFuncs []func()
	afterFuncs  []func()
}

func (e *Engine) AddAfterFunc(afterFunc func()) *Engine {
	e.afterFuncs = append(e.afterFuncs, afterFunc)
	return e
}

func (e *Engine) AddBeforeFunc(beforeFunc func()) *Engine {
	e.beforeFuncs = append(e.beforeFuncs, beforeFunc)
	return e
}

func (e *Engine) init() *Engine {
	serverConf := runFlags.GetServer()
	envOpts := []hello.Option{
		hello.WithVersion(runFlags.Version),
		hello.WithID(runFlags.Hostname),
		hello.WithEnv(runFlags.Environment.String()),
		hello.WithMetadata(serverConf.GetMetadata()),
	}
	if strings.EqualFold(serverConf.GetUseRandomID(), "true") {
		envOpts = append(envOpts, hello.WithID(strutil.RandomID()))
	}
	hello.SetEnvWithOption(envOpts...)
	return e
}

func (e *Engine) Start() {
	e.init()
	wg := new(sync.WaitGroup)
	for _, endpoint := range e.endpoints {
		endpoint.init()
	}
	for _, beforeFunc := range e.beforeFuncs {
		beforeFunc()
	}
	for _, endpoint := range e.endpoints {
		endpoint.start(wg)
	}
	wg.Wait()
	for _, afterFunc := range e.afterFuncs {
		afterFunc()
	}
}

func (e *endpoint) init() {
	e.helper = klog.NewHelper(klog.With(klog.GetLogger(),
		"service.name", e.serviceName,
		"service.id", hello.ID(),
		"caller", klog.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID()),
	)
	e.app, e.cleanup, e.err = e.wireAppFunc(e.serviceName, runFlags.Bootstrap, e.helper)
}

func (e *endpoint) start(wg *sync.WaitGroup) {
	if e.err != nil {
		e.helper.Errorw("msg", "endpoint init failed", "error", e.err)
		return
	}
	wg.Go(func() {
		defer e.cleanup()
		if err := e.app.Run(); err != nil {
			e.helper.Errorw("msg", "app run failed", "error", err)
		}
	})
}

func NewApp(serviceName string, d *data.Data, srvs server.Servers, bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, error) {
	opts := []kratos.Option{
		kratos.Name(serviceName),
		kratos.ID(hello.ID()),
		kratos.Version(hello.Version()),
		kratos.Metadata(hello.Metadata()),
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
	}

	if registry := d.Registry(); registry != nil {
		opts = append(opts, kratos.Registrar(registry))
	}

	for _, srv := range srvs {
		if httpSrv, ok := srv.(*http.Server); ok {
			server.BindSwagger(httpSrv, bc, helper)
			server.BindMetrics(httpSrv, bc, helper)
		}
	}

	return kratos.New(opts...), nil
}
