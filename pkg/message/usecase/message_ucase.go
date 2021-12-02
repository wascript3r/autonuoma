package usecase

import (
	"context"
	"html"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/ticket"
	"github.com/wascript3r/autonuoma/pkg/user"
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
		MessageInfo: &message.MessageInfo{
			User: &user.UserInfo{
				ID:        mf.UserMeta.ID,
				FirstName: mf.UserMeta.FirstName,
				LastName:  mf.UserMeta.LastName,
			},
			Content: mf.Content,
			Time:    mf.Time,
		},
	})
}

func (u *Usecase) Send(ctx context.Context, userID int, role domain.Role, req *message.SendReq) error {
	if role != domain.UserRole && role != domain.AgentRole {
		return domain.ErrInvalidUserRole
	}

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

	if (role == domain.UserRole && meta.ClientID != userID) || (role == domain.AgentRole && meta.AgentID != nil && *meta.AgentID != userID) {
		return ticket.TicketNotOwnedError
	}

	if meta.Status == domain.EndedTicketStatus {
		return ticket.TicketAlreadyEndedError
	} else if meta.Status == domain.CreatedTicketStatus {
		if role == domain.AgentRole {
			return ticket.TicketNotAcceptedError
		}
	} else if meta.Status != domain.AcceptedTicketStatus || meta.AgentID == nil {
		return domain.ErrInvalidTicketStatus
	}

	m := &domain.Message{
		TicketID: req.TicketID,
		UserID:   userID,
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
