package usecase

import (
	"context"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/trip"
)

type Usecase struct {
	tripRepo trip.Repository
}

func New(tr trip.Repository) *Usecase {
	return &Usecase{
		tripRepo: tr,
	}
}

func (u *Usecase) Start(ctx context.Context, req *trip.StartReq) (int, error) {
	var id, err = u.tripRepo.Start(ctx, req.StartTime, req.ReservationID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *Usecase) End(ctx context.Context, req *trip.EndReq) error {
	err := u.tripRepo.End(ctx, req.TripID, req.EndLat, req.EndLng)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) GetById(ctx context.Context, reservationID int) (*domain.Trip, error) {
	t := &domain.Trip{}
	t, err := u.tripRepo.GetByReservationId(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	return t, nil
}
