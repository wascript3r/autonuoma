package usecase

import (
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/room"
)

type Usecase struct {
	roomRepo room.Repository
}

func New(rr room.Repository) *Usecase {
	return &Usecase{
		roomRepo: rr,
	}
}

func (u *Usecase) Register(r domain.Room, c room.Config) error {
	if !domain.IsValidRoom(r) {
		return domain.ErrInvalidRoomName
	}

	return u.roomRepo.Set(r, c)
}

func (u *Usecase) GetName(r domain.Room) (string, error) {
	if !domain.IsValidRoom(r) {
		return "", domain.ErrInvalidRoomName
	}

	return u.roomRepo.GetName(r)
}
