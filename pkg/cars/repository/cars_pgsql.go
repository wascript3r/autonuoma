package repository

import (
	"context"
	"database/sql"
	"math/rand"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	getAllSQL    = "SELECT id, valstybiniai_numeriai, markė, modelis, pozicijos_platuma, pozicijos_ilguma FROM automobiliai WHERE pašalintas = false ORDER BY id ASC"
	getSingleSQL = "SELECT id, valstybiniai_numeriai, markė, modelis, spalva, pozicijos_platuma, pozicijos_ilguma, minutės_kaina, valandos_kaina, paros_kaina, kilometro_kaina, kondicionierius, usb, bluetooth, navigacija, vaikiška_kėdutė, pavarų_dėžė, kuro_tipas FROM automobiliai WHERE id = $1"
	removeCarSQL = "UPDATE automobiliai SET pašalintas = true WHERE id = $1"
	addCarSQL    = "INSERT INTO automobiliai (valstybiniai_numeriai, markė, modelis, spalva, minutės_kaina, valandos_kaina, paros_kaina, kilometro_kaina, kondicionierius, usb, bluetooth, navigacija, vaikiška_kėdutė, pavarų_dėžė, kuro_tipas, pašalintas, pozicijos_platuma, pozicijos_ilguma) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, false, $16, $17)"
)

type scanFunc = func(row pgsql.Row) (*domain.Car, error)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func scanRow(row pgsql.Row) (*domain.Car, error) {
	f := &domain.Car{}

	err := row.Scan(&f.ID, &f.LicensePlate, &f.Make, &f.Model, &f.Latitude, &f.Longitude)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return f, nil
}

func scanRows(rows *sql.Rows, scan scanFunc) ([]*domain.Car, error) {
	var fs []*domain.Car

	for rows.Next() {
		f, err := scan(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		fs = append(fs, f)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return fs, nil
}

func (p *PgRepo) GetAll(ctx context.Context) ([]*domain.Car, error) {
	rows, err := p.conn.QueryContext(ctx, getAllSQL)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, scanRow)
}

func (p PgRepo) GetSingle(ctx context.Context, carId int) (*domain.Car, error) {
	c := &domain.Car{}

	err := p.conn.QueryRowContext(ctx, getSingleSQL, carId).Scan(&c.ID, &c.LicensePlate, &c.Make, &c.Model, &c.Color, &c.Latitude, &c.Longitude, &c.MinutePrice, &c.HourPrice, &c.DayPrice, &c.KilometerPrice, &c.AirConditioning, &c.USB, &c.Bluetooth, &c.Navigation, &c.ChildSeat, &c.Fuel, &c.Gearbox)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return c, nil
}

func (p PgRepo) RemoveCar(ctx context.Context, carId int) (*domain.Car, error) {
	c := &domain.Car{}
	p.conn.ExecContext(ctx, removeCarSQL, carId)
	return c, nil
}

//"INSERT INTO automobiliai (valstybiniai_numeriai, markė, modelis, spalva, minutės_kaina, valandos_kaina, paros_kaina, kilometro_kaina, kondicionierius, usb, bluetooth, navigacija, vaikiška_kėdutė, pavarų_dėžė, kuro_tipas, pašalintas, pozicijos_platuma, pozicijos_ilguma) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, false, $16, $17)"

func (p PgRepo) AddCar(ctx context.Context, license_plate string, car_make string, car_model string, car_color string, minute_price float64, hour_price float64, day_price float64, kilometer_price float64, ac bool, usb bool, bluetooth bool, navigation bool, child_seat bool, gearbox int, fuel int) (*domain.Car, error) {
	c := &domain.Car{}

	positions := [5][2]float64{
		{54.9029023, 23.9598273},
		{54.9279151, 23.9610462},
		{54.938583, 23.8925454},
		{54.9148747, 23.8341097},
		{54.9450474, 23.8170732},
	}

	val := positions[rand.Intn(5)]
	x := val[0]
	y := val[1]

	_, err := p.conn.ExecContext(ctx, addCarSQL, license_plate, car_make, car_model, car_color, minute_price, hour_price, day_price, kilometer_price, ac, usb, bluetooth, navigation, child_seat, gearbox, fuel, x, y)

	if err != nil {
		print(err.Error())
		return nil, err
	}

	return c, nil
}
