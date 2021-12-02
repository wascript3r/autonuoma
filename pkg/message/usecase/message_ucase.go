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

	messageEventBus message.EventBus
	validate        ticket.Validate
}

func New(mr message.Repository, tr ticket.Repository, t time.Duration, meb message.EventBus, v ticket.Validate) *Usecase {
	return &Usecase{
		messageRepo: mr,
		ticketRepo:  tr,
		ctxTimeout:  t,

		messageEventBus: meb,
		validate:        v,
	}
}

func (u *Usecase) publishNewMessage(ctx context.Context, ticketID int, mf *domain.MessageFull) {
	u.messageEventBus.Publish(message.NewMessageEvent, ctx, &message.TicketMessage{
		TicketID: ticketID,
		MessageInfo: &ticket.MessageInfo{
			User: &ticket.UserInfo{
				ID:        mf.UserMeta.ID,
				FirstName: mf.UserMeta.FirstName,
				LastName:  mf.UserMeta.LastName,
			},
			Content: mf.Content,
			Time:    mf.Time,
		},
	})
}

func (u *Usecase) ClientSend(ctx context.Context, clientID int, req *message.ClientSendReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.messageRepo.NewTx(c)
	if err != nil {
		return err
	}

	tID, err := u.ticketRepo.GetLastActiveTicketIDTx(c, tx, clientID)
	if err != nil {
		if err == domain.ErrNotFound {
			return ticket.NoActiveTicketsError
		}
		return err
	}

	m := &domain.Message{
		TicketID: tID,
		UserID:   clientID,
		Content:  html.EscapeString(req.Message),
		Time:     time.Now(),
	}

	mf, err := u.messageRepo.InsertTx(c, tx, m)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	u.publishNewMessage(ctx, tID, mf)
	return nil
}

func (u *Usecase) AgentSend(ctx context.Context, agentID int, req *message.AgentSendReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.messageRepo.NewTx(c)
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
		return ticket.TicketNotAcceptedError
	} else if meta.Status != domain.AcceptedTicketStatus || meta.AgentID == nil {
		return domain.ErrInvalidTicketStatus
	} else if *meta.AgentID != agentID {
		return ticket.TicketNotOwnedError
	}

	m := &domain.Message{
		TicketID: req.TicketID,
		UserID:   agentID,
		Content:  html.EscapeString(req.Message),
		Time:     time.Now(),
	}

	mf, err := u.messageRepo.InsertTx(c, tx, m)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	u.publishNewMessage(ctx, req.TicketID, mf)
	return nil
}
