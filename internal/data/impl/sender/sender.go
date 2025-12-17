package sender

import (
	"context"

	"github.com/aide-family/magicbox/pointer"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type Sender interface {
	SendMessage(ctx context.Context, req *apiv1.JobSendMessageRequest) (*apiv1.JobSendReply, error)
	Close() error
}

// NewClusterSender creates a new cluster sender with GRPC client
func NewClusterSender(conn *grpc.ClientConn) Sender {
	return &clusterSender{
		conn: conn,
	}
}

// NewClusterHTTPSender creates a new cluster sender with HTTP client
func NewClusterHTTPSender(client *http.Client) Sender {
	return &clusterHTTPSender{
		client: client,
	}
}

type clusterSender struct {
	conn *grpc.ClientConn
}

func (c *clusterSender) SendMessage(ctx context.Context, req *apiv1.JobSendMessageRequest) (*apiv1.JobSendReply, error) {
	return apiv1.NewJobClient(c.conn).SendMessage(ctx, req)
}

func (c *clusterSender) Close() error {
	if pointer.IsNotNil(c.conn) {
		return c.conn.Close()
	}
	return nil
}

type clusterHTTPSender struct {
	client *http.Client
}

func (c *clusterHTTPSender) SendMessage(ctx context.Context, req *apiv1.JobSendMessageRequest) (*apiv1.JobSendReply, error) {
	return apiv1.NewJobHTTPClient(c.client).SendMessage(ctx, req)
}

func (c *clusterHTTPSender) Close() error {
	if pointer.IsNotNil(c.client) {
		return c.client.Close()
	}
	return nil
}
