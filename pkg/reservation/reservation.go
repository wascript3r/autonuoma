package reservation

type CreateReq struct {
	//StartAddress	string		`json:"startAddress" validate:"required"`
	//EndAddress		string		`json:"endAddress" validate:"required"`
	CarID int `json:"carID" validate:"required"`
}

type CancelReq struct {
	ReservationID int `json:"reservationID" validate:"required"`
}
