package domain

import (
	"errors"
	"time"
)

type TicketStatus int8

const (
	CreatedTicketStatus TicketStatus = iota
	AcceptedTicketStatus
	EndedTicketStatus
)

var ErrInvalidTicketStatus = errors.New("invalid ticket status")

func IsValidTicketStatus(ts TicketStatus) bool {
	switch ts {
	case CreatedTicketStatus, AcceptedTicketStatus, EndedTicketStatus:
		return true
	}
	return false
}

type Ticket struct {
	ID       int
	ClientID int
	AgentID  *int
	Created  time.Time
	Ended    *time.Time
}

type TicketMeta struct {
	Status   TicketStatus
	ClientID int
	AgentID  *int
	Ended    *time.Time
}

type TicketFull struct {
	ID           int
	Status       TicketStatus
	ClientMeta   *UserMeta
	FirstMessage string
	Time         time.Time
}
