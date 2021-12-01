package ticket

import (
	"context"
)

type Usecase interface {
	Create(ctx context.Context, clientID int, req *CreateReq) (int, error)
	Accept(ctx context.Context, agentID int, req *AcceptReq) error
}
