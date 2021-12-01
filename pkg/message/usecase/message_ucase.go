package usecase

import (
	"context"
	"html"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/ticket"
)

type Usecase struct {
	messageRepo message.Repository
	ticketRepo  ticket.Repository
	ctxTimeout  time.Duration

	validate ticket.Validate
}

func New(mr message.Repository, tr ticket.Repository, t time.Duration, v ticket.Validate) *Usecase {
	return &Usecase{
		messageRepo: mr,
		ticketRepo:  tr,
		ctxTimeout:  t,

		validate: v,
	}
}

func (u *Usecase) ClientSend(ctx context.Context, clientID int, req *message.ClientSendReq) (*domain.Message, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.messageRepo.NewTx(c)
	if err != nil {
		return nil, err
	}

	tID, err := u.ticketRepo.GetLastActiveTicketIDTx(c, tx, clientID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, ticket.TicketAlreadyEndedError
		}
		return nil, err
	}

	m := &domain.Message{
		TicketID: tID,
		UserID:   clientID,
		Content:  html.EscapeString(req.Message),
		Time:     time.Time{},
	}

	err = u.messageRepo.InsertTx(c, tx, m)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return m, nil
}
