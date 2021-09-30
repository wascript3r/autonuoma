package domain

import "time"

type Message struct {
	ID       int
	TicketID int
	Sender   *UserMeta
	Content  string
	Time     time.Time
}
