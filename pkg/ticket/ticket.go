package ticket

import (
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/review"
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

type EndReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

// GetFull

type TicketInfo struct {
	ID     int                 `json:"id"`
	Status domain.TicketStatus `json:"status"`
	Review *review.ReviewInfo  `json:"review"`
}

type GetFullReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

type GetFullRes struct {
	Ticket   *TicketInfo            `json:"ticket"`
	Messages []*message.MessageInfo `json:"messages"`
}

// GetAll

type TicketListInfo struct {
	ID           int                 `json:"id"`
	Status       domain.TicketStatus `json:"status"`
	Client       *user.UserInfo      `json:"client"`
	FirstMessage string              `json:"firstMessage"`
	Time         time.Time           `json:"time"`
}

type GetAllRes struct {
	Tickets []*TicketListInfo `json:"tickets"`
}
