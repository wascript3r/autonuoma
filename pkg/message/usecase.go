package message

import (
	"context"
)

type Usecase interface {
	ClientSend(ctx context.Context, clientID int, req *ClientSendReq) error
	AgentSend(ctx context.Context, agentID int, req *AgentSendReq) error
}
