package trip

import (
	"context"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"time"
)

type Repository interface {
	Start(ctx context.Context, endLng string, endLat string, reservationID int) (int, time.Time, error)
	End(ctx context.Context, tripID int, price float32) error
	GetByReservationId(ctx context.Context, reservationID int) (*domain.Trip, error)
}
