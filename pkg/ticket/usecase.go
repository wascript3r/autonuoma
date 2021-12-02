package ticket

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Create(ctx context.Context, clientID int, req *CreateReq) (int, error)
	Accept(ctx context.Context, agentID int, req *AcceptReq) error
	End(ctx context.Context, userID int, role domain.Role, req *EndReq) error
	GetMessages(ctx context.Context, userID int, role domain.Role, req *GetMessagesReq) (*GetMessagesRes, error)
	GetTickets(ctx context.Context, userID int, role domain.Role) (*GetTicketsRes, error)
}
