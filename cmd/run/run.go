// Package run is the run command for the Rabbit service
package run

import (
	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/server"
)

func NewCmd(defaultServerConfigBytes []byte) *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Rabbit service",
		Long: `启动 Rabbit 消息服务，提供统一的消息发送和管理能力。

Rabbit 是一个基于 Kratos 框架构建的分布式消息服务平台，支持多种消息通道
（邮件、Webhook、短信等）的统一管理和发送。通过命名空间（Namespace）实现
多租户隔离，支持配置文件和数据库两种存储模式，满足不同场景的部署需求。

主要功能：
  • 多通道消息发送：支持邮件、Webhook、短信等多种消息通道的统一管理
  • 模板化发送：支持消息模板配置，实现消息内容的动态渲染和复用
  • 异步消息处理：基于消息队列实现异步发送，提升系统吞吐量和可靠性
  • 配置管理：支持邮件服务器、Webhook 端点等通道配置的集中管理
  • 多租户隔离：通过命名空间实现不同业务或租户的配置和数据隔离

使用场景：
  • 企业级通知系统：统一管理各类业务通知（订单、告警、系统消息等）
  • 微服务消息中心：为微服务架构提供统一的消息发送能力
  • 多渠道推送平台：集成多种消息通道，实现消息的统一发送和管理

启动服务后，Rabbit 将监听配置的端口，提供 HTTP/gRPC API 接口供客户端调用。`,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runServer,
	}
	var bc conf.Bootstrap
	c := config.New(config.WithSource(
		env.NewSource(),
		conf.NewBytesSource(defaultServerConfigBytes),
	))
	if err := c.Load(); err != nil {
		flags.Helper.Errorw("msg", "load config failed", "error", err)
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		flags.Helper.Errorw("msg", "scan config failed", "error", err)
		panic(err)
	}

	flags.addFlags(runCmd, &bc)
	return runCmd
}

func runServer(_ *cobra.Command, _ []string) {
	flags.GlobalFlags = cmd.GetGlobalFlags()
	flags.applyToBootstrap()
	var bc conf.Bootstrap
	if strutil.IsNotEmpty(flags.configPath) {
		c := config.New(config.WithSource(
			env.NewSource(),
			file.NewSource(flags.configPath),
		))
		if err := c.Load(); err != nil {
			flags.Helper.Errorw("msg", "load config failed", "error", err)
			return
		}

		if err := c.Scan(&bc); err != nil {
			flags.Helper.Errorw("msg", "scan config failed", "error", err)
			return
		}
		flags.Bootstrap = &bc
	}

	serverConf := flags.GetServer()
	envOpts := []hello.Option{
		hello.WithVersion(flags.Version),
		hello.WithID(flags.Hostname),
		hello.WithName(serverConf.GetName()),
		hello.WithEnv(flags.Environment.String()),
		hello.WithMetadata(serverConf.GetMetadata()),
	}
	if serverConf.GetUseRandomID() == "true" {
		envOpts = append(envOpts, hello.WithID(strutil.RandomID()))
	}
	hello.SetEnvWithOption(envOpts...)

	helper := klog.NewHelper(klog.With(flags.Helper.Logger(),
		"cmd", "run",
		"service.name", hello.Name(),
		"service.id", hello.ID(),
		"caller", klog.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID()),
	)

	app, cleanup, err := wireApp(flags.Bootstrap, helper)
	if err != nil {
		flags.Helper.Errorw("msg", "wireApp failed", "error", err)
		return
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		flags.Helper.Errorw("msg", "app run failed", "error", err)
		return
	}
}

func newApp(d *data.Data, srvs server.Servers, bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, error) {
	defer hello.Hello()
	opts := []kratos.Option{
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
		kratos.Version(hello.Version()),
		kratos.ID(hello.ID()),
		kratos.Name(hello.Name()),
		kratos.Metadata(hello.Metadata()),
	}

	if registry := d.Registry(); registry != nil {
		opts = append(opts, kratos.Registrar(registry))
	}

	srvs.BindSwagger(bc, helper)
	srvs.BindMetrics(bc, helper)

	// 生成客户端配置
	if err := generateClientConfig(bc, srvs, helper); err != nil {
		helper.Warnw("msg", "generate client config failed", "error", err)
	}

	return kratos.New(opts...), nil
}
