package domain

import "time"

type FuelType int8

const (
	GasolineFuel FuelType = iota + 1
	DieselFuel
	ElectricityFuel
)

type GearboxType int8

const (
	AutomaticGearbox FuelType = iota + 1
	ManualGearbox
)

type Car struct {
	ID              int
	LicensePlate    string
	Make            string
	Model           string
	Color           string
	Latitude        float64
	Longitude       float64
	MinutePrice     float64
	HourPrice       float64
	DayPrice        float64
	KilometerPrice  float64
	AirConditioning bool
	USB             bool
	Bluetooth       bool
	Navigation      bool
	ChildSeat       bool
	Fuel            FuelType
	Gearbox         GearboxType
}

type CarTrip struct {
	FirstName string
	LastName  string
	Duration  time.Time
	Price     int
}

type CarStatistics struct {
	CarName string
}
