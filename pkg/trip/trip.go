package trip

import "time"

type StartReq struct {
	StartTime     time.Time `json:"startTime" validate:"required"`
	ReservationID int       `json:"reservationID" validate:"required"`
}

type StartRes struct {
	TripID int `json:"tripID"`
}

type EndReq struct {
	TripID int    `json:"tripID" validate:"required"`
	EndLng string `json:"endLng" validate:"required"`
	EndLat string `json:"endLat" validate:"required"`
}

type GetReq struct {
	ReservationID int `json:"reservationID" validate:"required"`
}
