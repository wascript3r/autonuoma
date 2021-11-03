package domain

import "time"

type Ticket struct {
	ID       int
	ClientID int
	AgentID  *int
	Created  time.Time
	Ended    *time.Time
}
