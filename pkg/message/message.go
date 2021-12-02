package message

import "github.com/wascript3r/autonuoma/pkg/ticket"

// Send

type ClientSendReq struct {
	Message string `json:"message" validate:"required,m_message"`
}

type AgentSendReq struct {
	TicketID int    `json:"ticketID" validate:"required"`
	Message  string `json:"message" validate:"required,m_message"`
}

type TicketMessage struct {
	TicketID int `json:"ticketID"`
	*ticket.MessageInfo
}
