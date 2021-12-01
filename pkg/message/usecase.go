package message

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	ClientSend(ctx context.Context, clientID int, req *ClientSendReq) (*domain.Message, error)
	AgentSend(ctx context.Context, agentID int, req *AgentSendReq) (*domain.Message, error)
}
