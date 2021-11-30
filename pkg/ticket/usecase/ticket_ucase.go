package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/ticket"
)

type Usecase struct {
	ticketRepo  ticket.Repository
	messageRepo message.Repository
	ctxTimeout  time.Duration

	validate ticket.Validate
}

func New(tr ticket.Repository, mr message.Repository, t time.Duration, v ticket.Validate) *Usecase {
	return &Usecase{
		ticketRepo:  tr,
		messageRepo: mr,
		ctxTimeout:  t,

		validate: v,
	}
}

func (u *Usecase) Create(ctx context.Context, clientID int, req *ticket.CreateReq) (int, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return 0, ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.ticketRepo.NewTx(c)
	if err != nil {
		return 0, err
	}

	_, err = u.ticketRepo.GetLastActiveTicketIDTx(c, tx, clientID)
	if err != domain.ErrNotFound {
		if err != nil {
			return 0, err
		}
		return 0, ticket.TicketStillActiveError
	}

	t := &domain.Ticket{
		ClientID: clientID,
		AgentID:  nil,
		Created:  time.Now(),
		Ended:    nil,
	}

	err = u.ticketRepo.InsertTx(c, tx, t)
	if err != nil {
		return 0, err
	}

	m := &domain.Message{
		TicketID: t.ID,
		UserID:   clientID,
		Content:  req.Message,
		Time:     time.Time{},
	}

	err = u.messageRepo.InsertTx(c, tx, m)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return t.ID, nil
}
