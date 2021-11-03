package domain

import "time"

type Message struct {
	ID       int
	TicketID int
	UserID   int
	Content  string
	Time     time.Time
}
