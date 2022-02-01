package trip

import (
	"context"
	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Start(ctx context.Context, req *StartReq) (*StartRes, error)
	End(ctx context.Context, req *EndReq) error
	GetById(ctx context.Context, reservationID int) (*domain.Trip, error)
}
