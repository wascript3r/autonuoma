package usecase

import (
	"context"
	"fmt"
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

func (u *Usecase) Start(ctx context.Context, req *trip.StartReq) (*trip.StartRes, error) {
	var id, createdAt, err = u.tripRepo.Start(ctx, req.EndLng, req.EndLat, req.ReservationID)
	fmt.Println("trip start")
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	return &trip.StartRes{TripID: id, CreatedAt: createdAt}, nil
}

func (u *Usecase) End(ctx context.Context, req *trip.EndReq) error {
	err := u.tripRepo.End(ctx, req.TripID, req.Price)
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
