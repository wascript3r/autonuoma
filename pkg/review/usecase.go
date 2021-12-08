package review

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Submit(ctx context.Context, userID int, role domain.Role, req *CreateReq) error
}
