package ticket

// Create

type CreateReq struct {
	Message string `json:"message" validate:"required,m_message"`
}

// Accept

type AcceptReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}

// EndAgent

type EndAgentReq struct {
	TicketID int `json:"ticketID" validate:"required"`
}
