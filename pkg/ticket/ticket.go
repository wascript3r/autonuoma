package ticket

import (
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/user"
)

// Create

type CreateReq struct {
	Message string `json:"message" validate:"required,m_message"`
}

// Accept

type AcceptReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

// End

type AgentEndReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

// GetMessages

type TicketStatus struct {
	ID    int  `json:"id"`
	Ended bool `json:"ended"`
}

type AgentGetMessagesReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

type GetMessagesRes struct {
	Ticket   *TicketStatus          `json:"ticket"`
	Messages []*message.MessageInfo `json:"messages"`
}

// GetTickets

type TicketInfo struct {
	ID           int                 `json:"id"`
	Status       domain.TicketStatus `json:"status"`
	Client       *user.UserInfo      `json:"client"`
	FirstMessage string              `json:"firstMessage"`
	Time         time.Time           `json:"time"`
}

type GetTicketsRes struct {
	Tickets []*TicketInfo `json:"tickets"`
}
