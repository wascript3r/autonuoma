package domain

import "time"

type UserTrip struct {
	ID    int
	Begin time.Time
	End   time.Time
	From  string
	To    string
	Price float32
}

type Trip struct {
	ID            int
	Begin         time.Time
	End           time.Time
	EndLng        string
	EndLat        string
	Duration      time.Time
	Price         float32
	ReservationID int
}
