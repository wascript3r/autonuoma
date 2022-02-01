package trip

import "time"

type StartReq struct {
	ReservationID int    `json:"reservationID" validate:"required"`
	EndLat        string `json:"endLat" validate:"required"`
	EndLng        string `json:"endLng" validate:"required"`
}

type StartRes struct {
	TripID    int       `json:"tripID"`
	CreatedAt time.Time `json:"createdAt"`
}

type EndReq struct {
	TripID int     `json:"tripID" validate:"required"`
	Price  float32 `json:"price" validate:"required"`
}

type GetReq struct {
	ReservationID int `json:"reservationID" validate:"required"`
}
