package ticket

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Create(ctx context.Context, clientID int, req *CreateReq) (int, error)
	Accept(ctx context.Context, agentID int, req *AcceptReq) error
	End(ctx context.Context, userID int, role domain.Role, req *EndReq) error
	GetFull(ctx context.Context, userID int, role domain.Role, req *GetFullReq) (*GetFullRes, error)
	GetAll(ctx context.Context, userID int, role domain.Role) (*GetAllRes, error)
}
