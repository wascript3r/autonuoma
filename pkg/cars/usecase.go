package cars

import (
	"context"
)

type Usecase interface {
	GetAll(ctx context.Context) (*GetAllRes, error)
	GetSingle(ctx context.Context, carId int) (*SingleCarRes, error)
	RemoveCar(ctx context.Context, carId int) (*SingleCarRes, error)
	AddCar(ctx context.Context, req *AddCarReq) (*SingleCarRes, error)
}
