package domain

import "time"

type TicketStatus int8

const (
	CreatedTicketStatus TicketStatus = iota
	AcceptedTicketStatus
	EndedTicketStatus
)

type Ticket struct {
	ID       int
	ClientID int
	AgentID  *int
	Created  time.Time
	Ended    *time.Time
}

type TicketMeta struct {
	Status  TicketStatus
	AgentID *int
	Ended   *time.Time
}
