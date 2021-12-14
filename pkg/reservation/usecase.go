package reservation

import (
	"context"
	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Create(ctx context.Context, req *CreateReq, uid int) (int, error)
	Cancel(ctx context.Context, reservationId int) error
	GetCurrent(ctx context.Context, userID int) (*domain.Reservation, error)
}
