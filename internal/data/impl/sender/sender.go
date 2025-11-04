package sender

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/merr"
)

type Sender interface {
	SendMessage(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error)
	Close() error
}

func NewClusterSender(conn *grpc.ClientConn, http *http.Client, protocol config.ClusterConfig_Protocol) Sender {
	return &clusterSender{
		conn:     conn,
		http:     http,
		protocol: protocol,
	}
}

type clusterSender struct {
	conn     *grpc.ClientConn
	http     *http.Client
	protocol config.ClusterConfig_Protocol
}

func (c *clusterSender) SendMessage(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error) {
	switch c.protocol {
	case config.ClusterConfig_GRPC:
		return apiv1.NewSenderClient(c.conn).SendMessage(ctx, req)
	case config.ClusterConfig_HTTP:
		return apiv1.NewSenderHTTPClient(c.http).SendMessage(ctx, req)
	default:
		return nil, merr.ErrorInternalServer("cluster sender unknown protocol %s", c.protocol)
	}
}

func (c *clusterSender) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	if c.http != nil {
		return c.http.Close()
	}
	return nil
}
