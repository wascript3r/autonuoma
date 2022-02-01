package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/cars"
)

type Usecase struct {
	carsRepo   cars.Repository
	ctxTimeout time.Duration
	validate   cars.Validate
}

func New(fr cars.Repository, t time.Duration, v cars.Validate) *Usecase {
	return &Usecase{
		carsRepo:   fr,
		ctxTimeout: t,
		validate:   v,
	}
}

func (u *Usecase) GetAll(ctx context.Context) (*cars.GetAllRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	fs, err := u.carsRepo.GetAll(c)
	if err != nil {
		return nil, err
	}

	carslist := make([]*cars.CarsListInfo, len(fs))
	for i, f := range fs {
		carslist[i] = &cars.CarsListInfo{
			ID:           f.ID,
			LicensePlate: f.LicensePlate,
			Make:         f.Make,
			Model:        f.Model,
			Latitude:     f.Latitude,
			Longitude:    f.Longitude,
		}
	}

	return &cars.GetAllRes{
		Cars: carslist,
	}, nil
}

func (u *Usecase) GetSingle(ctx context.Context, carId int) (*cars.SingleCarRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	fs, err := u.carsRepo.GetSingle(c, carId)
	if err != nil {
		return nil, err
	}

	isReserved := u.carsRepo.GetCarReservation(c, carId)

	return &cars.SingleCarRes{
		ID:              fs.ID,
		LicensePlate:    fs.LicensePlate,
		Make:            fs.Make,
		Model:           fs.Model,
		Color:           fs.Color,
		Latitude:        fs.Latitude,
		Longitude:       fs.Longitude,
		MinutePrice:     fs.MinutePrice,
		HourPrice:       fs.HourPrice,
		DayPrice:        fs.DayPrice,
		KilometerPrice:  fs.KilometerPrice,
		AirConditioning: fs.AirConditioning,
		USB:             fs.USB,
		Bluetooth:       fs.Bluetooth,
		Navigation:      fs.Navigation,
		ChildSeat:       fs.ChildSeat,
		Fuel:            fs.Fuel,
		Gearbox:         fs.Gearbox,
		IsReserved:      isReserved,
	}, nil
}

func (u *Usecase) RemoveCar(ctx context.Context, carId int) (*cars.SingleCarRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	_, err := u.carsRepo.RemoveCar(c, carId)
	if err != nil {
		print(err.Error())
		return nil, err
	}

	return nil, nil
}

func (u *Usecase) AddCar(ctx context.Context, req *cars.AddCarReq) (*cars.SingleCarRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	if err := u.validate.RawRequest(req); err != nil {
		return nil, cars.InvalidInputError
	}

	_, err := u.carsRepo.AddCar(c, req.LicensePlate, req.Make, req.Model, req.Color, req.MinutePrice, req.HourPrice, req.DayPrice, req.KilometerPrice, req.AirConditioning, req.USB, req.Bluetooth, req.Navigation, req.ChildSeat, req.Gearbox, req.Fuel)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *Usecase) UpdateCar(ctx context.Context, req *cars.UpdateCarReq) (*cars.SingleCarRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	if err := u.validate.RawRequest(req); err != nil {
		return nil, cars.InvalidInputError
	}

	_, err := u.carsRepo.UpdateCar(c, req.Id, req.LicensePlate, req.Make, req.Model, req.Color, req.MinutePrice, req.HourPrice, req.DayPrice, req.KilometerPrice, req.AirConditioning, req.USB, req.Bluetooth, req.Navigation, req.ChildSeat, req.Gearbox, req.Fuel)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *Usecase) CarTrips(ctx context.Context, carId int) (*cars.CarTripsRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	fs, err := u.carsRepo.CarTrips(c, carId)
	if err != nil {
		return nil, err
	}

	tripsList := make([]*cars.CarTripsInfo, len(fs))
	for i, f := range fs {
		tripsList[i] = &cars.CarTripsInfo{
			FirstName: f.FirstName,
			LastName:  f.LastName,
			Duration:  f.Duration.Format("15:04:05"),
			Price:     f.Price,
		}
	}

	return &cars.CarTripsRes{
		Trips: tripsList,
	}, nil
}

func (u *Usecase) Statistics(ctx context.Context) (*cars.StatisticsRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	fs, err := u.carsRepo.Statistics(c)
	if err != nil {
		return nil, err
	}

	statisticsList := make([]*cars.CarStatisticsInfo, len(fs))
	for i, f := range fs {
		statisticsList[i] = &cars.CarStatisticsInfo{
			CarName: f.CarName,
		}
	}

	return &cars.StatisticsRes{
		Statistics: statisticsList,
	}, nil
}
