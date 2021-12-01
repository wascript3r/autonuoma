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
	ticketRepo  ticket.Repository
	messageRepo message.Repository
	ctxTimeout  time.Duration

	ticketEventBus ticket.EventBus
	validate       ticket.Validate
}

func New(tr ticket.Repository, mr message.Repository, t time.Duration, teb ticket.EventBus, v ticket.Validate) *Usecase {
	return &Usecase{
		ticketRepo:  tr,
		messageRepo: mr,
		ctxTimeout:  t,

		ticketEventBus: teb,
		validate:       v,
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
		Content:  html.EscapeString(req.Message),
		Time:     time.Now(),
	}

	err = u.messageRepo.InsertTx(c, tx, m)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	u.ticketEventBus.Publish(ticket.NewTicketEvent, ctx)

	return t.ID, nil
}

func (u *Usecase) Accept(ctx context.Context, agentID int, req *ticket.AcceptReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.ticketRepo.NewTx(c)
	if err != nil {
		return err
	}

	meta, err := u.ticketRepo.GetTicketMetaTx(c, tx, req.TicketID)
	if err != nil {
		if err == domain.ErrNotFound {
			return ticket.TicketNotFoundError
		}
		return err
	}

	if meta.Status != domain.CreatedTicketStatus {
		if meta.Status == domain.AcceptedTicketStatus {
			return ticket.TicketAlreadyAcceptedError
		} else if meta.Status == domain.EndedTicketStatus {
			return ticket.TicketAlreadyEndedError
		}
		return ticket.UnknownError
	}

	err = u.ticketRepo.SetAgentTx(c, tx, req.TicketID, agentID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) EndClient(ctx context.Context, clientID int) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.ticketRepo.NewTx(c)
	if err != nil {
		return err
	}

	id, err := u.ticketRepo.GetLastActiveTicketIDTx(c, tx, clientID)
	if err != nil {
		if err == domain.ErrNotFound {
			return ticket.NoActiveTicketsError
		}
		return err
	}

	err = u.ticketRepo.SetEndedTx(c, tx, id, time.Now())
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) EndAgent(ctx context.Context, agentID int, req *ticket.EndAgentReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.ticketRepo.NewTx(c)
	if err != nil {
		return err
	}

	meta, err := u.ticketRepo.GetTicketMetaTx(c, tx, req.TicketID)
	if err != nil {
		if err == domain.ErrNotFound {
			return ticket.TicketNotFoundError
		}
		return err
	}

	if meta.Status == domain.EndedTicketStatus {
		return ticket.TicketAlreadyEndedError
	} else if meta.Status == domain.CreatedTicketStatus {
		err = u.ticketRepo.SetAgentEndedTx(c, tx, req.TicketID, agentID, time.Now())
	} else if meta.Status == domain.AcceptedTicketStatus && meta.AgentID != nil {
		if *meta.AgentID != agentID {
			return ticket.TicketNotOwnedError
		}
		err = u.ticketRepo.SetEndedTx(c, tx, req.TicketID, time.Now())
	} else {
		return ticket.UnknownError
	}

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
