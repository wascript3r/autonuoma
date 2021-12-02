package message

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Send(ctx context.Context, userID int, role domain.Role, req *SendReq) error
}
