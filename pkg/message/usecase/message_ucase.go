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
			return nil, ticket.NoActiveTicketsError
		}
		return nil, err
	}

	m := &domain.Message{
		TicketID: tID,
		UserID:   clientID,
		Content:  html.EscapeString(req.Message),
		Time:     time.Now(),
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

func (u *Usecase) AgentSend(ctx context.Context, agentID int, req *message.AgentSendReq) (*domain.Message, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.messageRepo.NewTx(c)
	if err != nil {
		return nil, err
	}

	meta, err := u.ticketRepo.GetTicketMetaTx(c, tx, req.TicketID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, ticket.TicketNotFoundError
		}
		return nil, err
	}

	if meta.Status == domain.EndedTicketStatus {
		return nil, ticket.TicketAlreadyEndedError
	} else if meta.Status == domain.CreatedTicketStatus {
		return nil, ticket.TicketNotAcceptedError
	} else if meta.Status != domain.AcceptedTicketStatus || meta.AgentID == nil {
		return nil, ticket.UnknownError
	} else if *meta.AgentID != agentID {
		return nil, ticket.TicketNotOwnedError
	}

	m := &domain.Message{
		TicketID: req.TicketID,
		UserID:   agentID,
		Content:  html.EscapeString(req.Message),
		Time:     time.Now(),
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
