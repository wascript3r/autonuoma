package message

// Send

type ClientSendReq struct {
	Message string `json:"message" validate:"required,m_message"`
}

type AgentSendReq struct {
	TicketID int    `json:"ticketID" validate:"required"`
	Message  string `json:"message" validate:"required,m_message"`
}
