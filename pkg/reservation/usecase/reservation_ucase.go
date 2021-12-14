package usecase

import (
	"context"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/reservation"
)

type Usecase struct {
	resRepo reservation.Repository
}

func New(rr reservation.Repository) *Usecase {
	return &Usecase{
		resRepo: rr,
	}
}

func (u *Usecase) Create(ctx context.Context, req *reservation.CreateReq, uid int) (int, error) {
	var id, err = u.resRepo.Create(ctx, req.CarID, uid)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *Usecase) Cancel(ctx context.Context, reservationID int) error {
	err := u.resRepo.Cancel(ctx, reservationID)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) GetCurrent(ctx context.Context, userID int) (*domain.Reservation, error) {
	r := &domain.Reservation{}
	r, err := u.resRepo.GetCurrent(ctx, userID)
	if err != nil {
		return nil, err
	}
	return r, nil
}
