package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/ticket"
	"github.com/wascript3r/autonuoma/pkg/user"
)

type Usecase struct {
	ticketRepo ticket.Repository
	userRepo   user.Repository
	ctxTimeout time.Duration

	validate ticket.Validate
}

func New(tr ticket.Repository, ur user.Repository, t time.Duration, v ticket.Validate) *Usecase {
	return &Usecase{
		ticketRepo: tr,
		userRepo:   ur,
		ctxTimeout: t,

		validate: v,
	}
}

func (u *Usecase) Create(ctx context.Context, userID int, req *ticket.CreateReq) (int, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return 0, ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.userRepo.NewTx(c)
	if err != nil {
		return 0, err
	}

	ended, err := u.ticketRepo.IsCurrTicketEndedTx(c, tx, userID)
	if err != domain.ErrNotFound {
		if err != nil {
			return 0, err
		}

		if !ended {
			return 0, user.TicketStillActiveError
		}
	}

	t := &domain.Ticket{
		ClientID: userID,
		AgentID:  nil,
		Created:  time.Now(),
		Ended:    nil,
	}

	err = u.ticketRepo.InsertTx(ctx, tx, t)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return t.ID, nil
}
