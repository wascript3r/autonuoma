package cars

import "github.com/wascript3r/autonuoma/pkg/domain"

type SingleCarReq struct {
	CarID int `json:"id" validate:"required"`
}

// GetAll

type CarsListInfo struct {
	ID           int    `json:"id"`
	LicensePlate string `json:"license_plate"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Latitude     string `json:"lat"`
	Longitude    string `json:"lng"`
}

type GetAllRes struct {
	Cars []*CarsListInfo `json:"cars"`
}

// GetSingle

type SingleCarRes struct {
	ID              int
	LicensePlate    string
	Make            string
	Model           string
	Color           string
	Latitude        string
	Longitude       string
	MinutePrice     float64
	HourPrice       float64
	DayPrice        float64
	KilometerPrice  float64
	AirConditioning bool
	USB             bool
	Bluetooth       bool
	Navigation      bool
	ChildSeat       bool
	Fuel            domain.FuelType
	Gearbox         domain.GearboxType
}

// Add car
type AddCarReq struct {
	LicensePlate    string  `json:"license_plate" validate:"required"`
	Make            string  `json:"make" validate:"required"`
	Model           string  `json:"model" validate:"required"`
	Color           string  `json:"color" validate:"required"`
	MinutePrice     float64 `json:"minute_price" validate:"required"`
	HourPrice       float64 `json:"hour_price" validate:"required"`
	DayPrice        float64 `json:"day_price" validate:"required"`
	KilometerPrice  float64 `json:"kilometer_price" validate:"required"`
	AirConditioning bool    `json:"air_conditioning"`
	USB             bool    `json:"usb"`
	Bluetooth       bool    `json:"bluetooth"`
	Navigation      bool    `json:"navigation"`
	ChildSeat       bool    `json:"child_seat"`
	Fuel            int     `json:"fuel" validate:"required"`
	Gearbox         int     `json:"gearbox" validate:"required"`
}

// Update car
type UpdateCarReq struct {
	Id              int     `json:"id" validate:"required"`
	LicensePlate    string  `json:"license_plate" validate:"required"`
	Make            string  `json:"make" validate:"required"`
	Model           string  `json:"model" validate:"required"`
	Color           string  `json:"color" validate:"required"`
	MinutePrice     float64 `json:"minute_price" validate:"required"`
	HourPrice       float64 `json:"hour_price" validate:"required"`
	DayPrice        float64 `json:"day_price" validate:"required"`
	KilometerPrice  float64 `json:"kilometer_price" validate:"required"`
	AirConditioning bool    `json:"air_conditioning"`
	USB             bool    `json:"usb"`
	Bluetooth       bool    `json:"bluetooth"`
	Navigation      bool    `json:"navigation"`
	ChildSeat       bool    `json:"child_seat"`
	Fuel            int     `json:"fuel" validate:"required"`
	Gearbox         int     `json:"gearbox" validate:"required"`
}

// Car trips
type CarTripsInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Duration  string `json:"duration"`
	Price     int    `json:"price"`
}

type CarTripsRes struct {
	Trips []*CarTripsInfo `json:"trips"`
}

// Car statistics
type CarStatisticsInfo struct {
	CarName string `json:"car"`
}

type StatisticsRes struct {
	Statistics []*CarStatisticsInfo `json:"statistics"`
}
