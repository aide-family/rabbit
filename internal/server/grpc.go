package server

import (
	"github.com/aide-family/rabbit/internal/conf"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bc *conf.Bootstrap, helper *klog.Helper) *grpc.Server {
	serverConf := bc.GetServer()
	grpcConf := serverConf.GetGrpc()

	opts := []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(helper.Logger()),
			tracing.Server(),
			metadata.Server(),
		),
	}
	if grpcConf.GetNetwork() != "" {
		opts = append(opts, grpc.Network(grpcConf.GetNetwork()))
	}
	if grpcConf.GetAddress() != "" {
		opts = append(opts, grpc.Address(grpcConf.GetAddress()))
	}
	if grpcConf.GetTimeout() != nil {
		opts = append(opts, grpc.Timeout(grpcConf.GetTimeout().AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	return srv
}
