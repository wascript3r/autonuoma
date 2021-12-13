package cars

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*domain.Car, error)
	GetSingle(ctx context.Context, carId int) (*domain.Car, error)
	RemoveCar(ctx context.Context, carId int) (*domain.Car, error)
	AddCar(ctx context.Context, license_plate string, car_make string, car_model string, car_color string, minute_price float64, hour_price float64, day_price float64, kilometer_price float64, ac bool, usb bool, bluetooth bool, navigation bool, child_seat bool, gearbox int, fuel int) (*domain.Car, error)
}
