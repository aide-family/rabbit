package sender

import (
	"context"

	"github.com/aide-family/magicbox/pointer"
	"google.golang.org/grpc"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type Sender interface {
	SendMessage(ctx context.Context, req *apiv1.JobSendMessageRequest) (*apiv1.JobSendReply, error)
	Close() error
}

func NewClusterSender(conn *grpc.ClientConn) Sender {
	return &clusterSender{
		conn: conn,
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
