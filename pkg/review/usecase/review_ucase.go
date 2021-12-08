package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/review"
	"github.com/wascript3r/autonuoma/pkg/ticket"
)

type Usecase struct {
	reviewRepo review.Repository
	ticketRepo ticket.Repository
	ctxTimeout time.Duration

	validate ticket.Validate
}

func New(rr review.Repository, tr ticket.Repository, t time.Duration, v ticket.Validate) *Usecase {
	return &Usecase{
		reviewRepo: rr,
		ticketRepo: tr,
		ctxTimeout: t,

		validate: v,
	}
}

func (u *Usecase) Submit(ctx context.Context, userID int, role domain.Role, req *review.CreateReq) error {
	if role != domain.ClientRole {
		return domain.ErrInvalidUserRole
	}

	if err := u.validate.RawRequest(req); err != nil {
		return ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.reviewRepo.NewTx(c)
	if err != nil {
		return err
	}

	meta, err := u.ticketRepo.GetMetaTx(c, tx, req.TicketID)
	if err != nil {
		if err == domain.ErrNotFound {
			return ticket.TicketNotFoundError
		}
		return err
	}

	if meta.ClientID != userID {
		return ticket.TicketNotOwnedError
	}

	if meta.Status != domain.EndedTicketStatus {
		return ticket.TicketNotEndedError
	}

	_, err = u.reviewRepo.GetByTicketTx(ctx, tx, req.TicketID)
	if err != domain.ErrNotFound {
		if err != nil {
			return err
		}
		return review.ReviewAlreadySubmittedError
	}

	r := &domain.Review{
		TicketID: req.TicketID,
		Stars:    req.Stars,
		Comment:  req.Comment,
		Time:     time.Now(),
	}

	err = u.reviewRepo.InsertTx(c, tx, r)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
