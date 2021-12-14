package reservation

import (
	"context"
	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Repository interface {
	Create(ctx context.Context, cardID int, userID int) (int, error)
	Cancel(ctx context.Context, reservationID int) error
	GetCurrent(ctx context.Context, userID int) (*domain.Reservation, error)
}
