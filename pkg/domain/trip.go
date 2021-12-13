package domain

import "time"

type Trip struct {
	ID    int
	Begin time.Time
	End   time.Time
	From  string
	To    string
	Price float32
}
