package ticket

import "time"

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

type UserInfo struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type MessageInfo struct {
	User    *UserInfo `json:"user"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type AgentGetMessagesReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

type GetMessagesRes struct {
	TicketEnded bool           `json:"ticketEnded"`
	Messages    []*MessageInfo `json:"messages"`
}
