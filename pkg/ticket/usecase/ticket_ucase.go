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

	_, err = u.messageRepo.InsertTx(c, tx, m)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	u.ticketEventBus.Publish(ticket.NewTicketEvent, ctx, t.ID)

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
		return domain.ErrInvalidTicketStatus
	}

	err = u.ticketRepo.SetAgentTx(c, tx, req.TicketID, agentID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	u.ticketEventBus.Publish(ticket.AcceptedTicketEvent, ctx, req.TicketID)
	return nil
}

func (u *Usecase) End(ctx context.Context, userID int, role domain.Role, req *ticket.EndReq) error {
	if role != domain.UserRole && role != domain.AgentRole {
		return domain.ErrInvalidUserRole
	}

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

	if (role == domain.UserRole && meta.ClientID != userID) || (role == domain.AgentRole && meta.AgentID != nil && *meta.AgentID != userID) {
		return ticket.TicketNotOwnedError
	}

	if meta.Status == domain.EndedTicketStatus {
		return ticket.TicketAlreadyEndedError
	} else if meta.Status == domain.AcceptedTicketStatus && meta.AgentID != nil {
		err = u.ticketRepo.SetEndedTx(c, tx, req.TicketID, time.Now())
	} else if meta.Status == domain.CreatedTicketStatus {
		if role == domain.AgentRole {
			err = u.ticketRepo.SetAgentEndedTx(c, tx, req.TicketID, userID, time.Now())
		} else {
			err = u.ticketRepo.SetEndedTx(c, tx, req.TicketID, time.Now())
		}
	} else {
		return domain.ErrInvalidTicketStatus
	}

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	u.ticketEventBus.Publish(ticket.EndedTicketEvent, ctx, req.TicketID)
	return nil
}

func (u *Usecase) GetMessages(ctx context.Context, userID int, role domain.Role, req *ticket.GetMessagesReq) (*ticket.GetMessagesRes, error) {
	if role != domain.UserRole && role != domain.AgentRole {
		return nil, domain.ErrInvalidUserRole
	}

	if err := u.validate.RawRequest(req); err != nil {
		return nil, ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	meta, err := u.ticketRepo.GetTicketMeta(c, req.TicketID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, ticket.TicketNotFoundError
		}
		return nil, err
	}

	if role == domain.UserRole && meta.ClientID != userID {
		return nil, ticket.TicketNotOwnedError
	} else if !domain.IsValidTicketStatus(meta.Status) {
		return nil, domain.ErrInvalidTicketStatus
	}

	ms, err := u.messageRepo.GetTicketMessages(ctx, req.TicketID)
	if err != nil {
		return nil, err
	}

	messages := make([]*message.MessageInfo, len(ms))
	for i, m := range ms {
		messages[i] = &message.MessageInfo{
			User: &user.UserInfo{
				ID:        m.UserMeta.ID,
				FirstName: m.UserMeta.FirstName,
				LastName:  m.UserMeta.LastName,
			},
			Content: m.Content,
			Time:    m.Time,
		}
	}

	return &ticket.GetMessagesRes{
		Ticket: &ticket.TicketStatus{
			ID:     req.TicketID,
			Status: meta.Status,
		},
		Messages: messages,
	}, nil
}

func (u *Usecase) GetTickets(ctx context.Context, userID int, role domain.Role) (*ticket.GetTicketsRes, error) {
	if role != domain.UserRole && role != domain.AgentRole {
		return nil, domain.ErrInvalidUserRole
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	var (
		ts  []*domain.TicketFull
		err error
	)

	if role == domain.UserRole {
		ts, err = u.ticketRepo.GetUserTickets(c, userID)
	} else {
		ts, err = u.ticketRepo.GetTickets(c)
	}

	if err != nil {
		return nil, err
	}

	tickets := make([]*ticket.TicketInfo, len(ts))
	for i, t := range ts {
		tickets[i] = &ticket.TicketInfo{
			ID:     t.ID,
			Status: t.Status,
			Client: &user.UserInfo{
				ID:        t.ClientMeta.ID,
				FirstName: t.ClientMeta.FirstName,
				LastName:  t.ClientMeta.LastName,
			},
			FirstMessage: t.FirstMessage,
			Time:         t.Time,
		}
	}

	return &ticket.GetTicketsRes{
		Tickets: tickets,
	}, nil
}
