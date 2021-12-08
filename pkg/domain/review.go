package domain

import "time"

type Review struct {
	ID       int
	TicketID int
	Stars    int
	Comment  *string
	Time     time.Time
}
