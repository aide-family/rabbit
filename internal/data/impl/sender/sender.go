package sender

import (
	"context"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type Sender interface {
	SendMessage(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error)
	Close() error
}
