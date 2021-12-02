package ticket

import "github.com/wascript3r/autonuoma/pkg/message"

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

type TicketInfo struct {
	ID    int  `json:"id"`
	Ended bool `json:"ended"`
}

type AgentGetMessagesReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

type GetMessagesRes struct {
	Ticket   *TicketInfo            `json:"ticket"`
	Messages []*message.MessageInfo `json:"messages"`
}
