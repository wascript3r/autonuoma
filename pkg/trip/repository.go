package trip

import (
	"context"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"time"
)

type Repository interface {
	Start(ctx context.Context, startTime time.Time, reservationID int) (int, error)
	End(ctx context.Context, tripID int, endLat string, endLng string) error
	GetByReservationId(ctx context.Context, reservationID int) (*domain.Trip, error)
}
