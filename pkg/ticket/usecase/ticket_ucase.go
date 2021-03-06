package usecase

import (
	"context"
	"html"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/review"
	"github.com/wascript3r/autonuoma/pkg/ticket"
	"github.com/wascript3r/autonuoma/pkg/user"
)

type Usecase struct {
	ticketRepo  ticket.Repository
	messageRepo message.Repository
	reviewRepo  review.Repository
	ctxTimeout  time.Duration

	ticketEventBus ticket.EventBus
	validate       ticket.Validate
}

func New(tr ticket.Repository, mr message.Repository, rr review.Repository, t time.Duration, teb ticket.EventBus, v ticket.Validate) *Usecase {
	return &Usecase{
		ticketRepo:  tr,
		messageRepo: mr,
		reviewRepo:  rr,
		ctxTimeout:  t,

		ticketEventBus: teb,
		validate:       v,
	}
}

func (u *Usecase) Create(ctx context.Context, clientID int, req *ticket.CreateReq) (*ticket.CreateRes, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.ticketRepo.NewTx(c)
	if err != nil {
		return nil, err
	}

	_, err = u.ticketRepo.GetLastActiveIDTx(c, tx, clientID)
	if err != domain.ErrNotFound {
		if err != nil {
			return nil, err
		}
		return nil, ticket.TicketStillActiveError
	}

	t := &domain.Ticket{
		ClientID: clientID,
		AgentID:  nil,
		Created:  time.Now(),
		Ended:    nil,
	}

	err = u.ticketRepo.InsertTx(c, tx, t)
	if err != nil {
		return nil, err
	}

	m := &domain.Message{
		TicketID: t.ID,
		UserID:   clientID,
		Content:  html.EscapeString(req.Message),
		Time:     time.Now(),
	}

	_, err = u.messageRepo.InsertTx(c, tx, m)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	u.ticketEventBus.Publish(ticket.NewTicketEvent, ctx, t.ID)

	return &ticket.CreateRes{
		TicketID: t.ID,
	}, nil
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

	meta, err := u.ticketRepo.GetMetaTx(c, tx, req.TicketID)
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
	if role != domain.ClientRole && role != domain.AgentRole {
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

	meta, err := u.ticketRepo.GetMetaTx(c, tx, req.TicketID)
	if err != nil {
		if err == domain.ErrNotFound {
			return ticket.TicketNotFoundError
		}
		return err
	}

	if (role == domain.ClientRole && meta.ClientID != userID) || (role == domain.AgentRole && meta.AgentID != nil && *meta.AgentID != userID) {
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

func (u *Usecase) GetFull(ctx context.Context, userID int, role domain.Role, req *ticket.GetFullReq) (*ticket.GetFullRes, error) {
	if role != domain.ClientRole && role != domain.AgentRole {
		return nil, domain.ErrInvalidUserRole
	}

	if err := u.validate.RawRequest(req); err != nil {
		return nil, ticket.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	meta, err := u.ticketRepo.GetMeta(c, req.TicketID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, ticket.TicketNotFoundError
		}
		return nil, err
	}

	if role == domain.ClientRole && meta.ClientID != userID {
		return nil, ticket.TicketNotOwnedError
	} else if !domain.IsValidTicketStatus(meta.Status) {
		return nil, domain.ErrInvalidTicketStatus
	}

	rs, err := u.reviewRepo.GetByTicket(c, req.TicketID)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}

	ms, err := u.messageRepo.GetByTicket(c, req.TicketID)
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

	res := &ticket.GetFullRes{
		Ticket: &ticket.TicketInfo{
			ID:      req.TicketID,
			Status:  meta.Status,
			AgentID: meta.AgentID,
			Review:  nil,
		},
		Messages: messages,
	}
	if rs != nil {
		res.Ticket.Review = &review.ReviewInfo{
			Stars:   rs.Stars,
			Comment: rs.Comment,
		}
	}

	return res, nil
}

func (u *Usecase) GetAll(ctx context.Context, userID int, role domain.Role) (*ticket.GetAllRes, error) {
	if role != domain.ClientRole && role != domain.AgentRole {
		return nil, domain.ErrInvalidUserRole
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	var (
		ts  []*domain.TicketFull
		err error
	)

	if role == domain.ClientRole {
		ts, err = u.ticketRepo.GetByUser(c, userID)
	} else {
		ts, err = u.ticketRepo.GetAll(c)
	}

	if err != nil {
		return nil, err
	}

	tickets := make([]*ticket.TicketListInfo, len(ts))
	for i, t := range ts {
		tickets[i] = &ticket.TicketListInfo{
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

	return &ticket.GetAllRes{
		Tickets: tickets,
	}, nil
}
