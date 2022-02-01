package reservation

type CreateReq struct {
	CarID int `json:"carID" validate:"required"`
}

type CancelReq struct {
	ReservationID int `json:"reservationID" validate:"required"`
}
