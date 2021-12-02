package ticket

import (
	"context"
)

type Usecase interface {
	Create(ctx context.Context, clientID int, req *CreateReq) (int, error)
	Accept(ctx context.Context, agentID int, req *AcceptReq) error
	ClientEnd(ctx context.Context, clientID int) error
	AgentEnd(ctx context.Context, agentID int, req *AgentEndReq) error
	ClientGetMessages(ctx context.Context, clientID int) (*GetMessagesRes, error)
	AgentGetMessages(ctx context.Context, req *AgentGetMessagesReq) (*GetMessagesRes, error)
	GetTickets(ctx context.Context) (*GetTicketsRes, error)
}
