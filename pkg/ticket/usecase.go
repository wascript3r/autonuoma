package ticket

import (
	"context"
)

type Usecase interface {
	Create(ctx context.Context, userID int, req *CreateReq) (int, error)
}
