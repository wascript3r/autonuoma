package domain

import "time"

type Reservation struct {
	ID           int
	CreatedAt    time.Time
	CanceledAt   time.Time
	StartAddress string
	EndAddress   string
	CarID        int
	UserID       int
}
