package repository

import (
	"context"
	"database/sql"
	"math/rand"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	getAllSQL     = "SELECT id, valstybiniai_numeriai, markė, modelis FROM automobiliai WHERE pašalintas = false ORDER BY id ASC"
	getSingleSQL  = "SELECT id, valstybiniai_numeriai, markė, modelis, spalva, pozicijos_platuma, pozicijos_ilguma, minutės_kaina, valandos_kaina, paros_kaina, kilometro_kaina, kondicionierius, usb, bluetooth, navigacija, vaikiška_kėdutė, pavarų_dėžė, kuro_tipas FROM automobiliai WHERE id = $1"
	removeCarSQL  = "UPDATE automobiliai SET pašalintas = true WHERE id = $1"
	addCarSQL     = "INSERT INTO automobiliai (valstybiniai_numeriai, markė, modelis, spalva, minutės_kaina, valandos_kaina, paros_kaina, kilometro_kaina, kondicionierius, usb, bluetooth, navigacija, vaikiška_kėdutė, pavarų_dėžė, kuro_tipas, pašalintas, pozicijos_platuma, pozicijos_ilguma) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, false, $16, $17)"
	updateCarSQL  = "UPDATE automobiliai SET valstybiniai_numeriai = $1, markė = $2, modelis = $3, spalva = $4, minutės_kaina = $5, valandos_kaina = $6, paros_kaina = $7, kilometro_kaina = $8, kondicionierius = $9, usb = $10, bluetooth = $11, navigacija = $12, vaikiška_kėdutė = $13, pavarų_dėžė = $14, kuro_tipas = $15 WHERE id = $16"
	carTripsSQL   = "SELECT v.vardas, v.pavardė, k.trukmė, k.kaina FROM kelionės k INNER JOIN rezervacijos r ON r.id = k.fk_rezervacija INNER JOIN vartotojai v ON v.id = r.fk_vartotojas WHERE r.fk_automobilis = $1"
	statisticsSQL = "SELECT x.pavadinimas FROM ((SELECT 1 nr, CONCAT(a.markė, ' ', a.modelis, ' (', a.valstybiniai_numeriai, ')') pavadinimas FROM rezervacijos r INNER JOIN automobiliai a ON a.id = r.fk_automobilis GROUP BY a.markė, a.modelis, a.valstybiniai_numeriai ORDER BY COUNT(*) DESC LIMIT 1) UNION (SELECT 2 nr, CONCAT(a.markė, ' ', a.modelis, ' (', a.valstybiniai_numeriai, ')') pavadinimas FROM rezervacijos r INNER JOIN automobiliai a ON a.id = r.fk_automobilis GROUP BY a.markė, a.modelis, a.valstybiniai_numeriai ORDER BY COUNT(*) ASC LIMIT 1) UNION (SELECT 3 nr, CONCAT(a.markė, ' ', a.modelis, ' (', a.valstybiniai_numeriai, ')') pavadinimas FROM rezervacijos r INNER JOIN automobiliai a ON a.id = r.fk_automobilis INNER JOIN kelionės k ON k.fk_rezervacija = r.id GROUP BY a.markė, a.modelis, a.valstybiniai_numeriai ORDER BY SUM(k.kaina) DESC LIMIT 1) UNION (SELECT 4 nr, CONCAT(a.markė, ' ', a.modelis, ' (', a.valstybiniai_numeriai, ')') pavadinimas FROM rezervacijos r INNER JOIN automobiliai a ON a.id = r.fk_automobilis INNER JOIN kelionės k ON k.fk_rezervacija = r.id GROUP BY a.markė, a.modelis, a.valstybiniai_numeriai ORDER BY SUM(k.kaina) ASC LIMIT 1)) x ORDER BY x.nr DESC"
)

type scanFunc = func(row pgsql.Row) (*domain.Car, error)
type scanFunc2 = func(row pgsql.Row) (*domain.CarTrip, error)
type scanFunc3 = func(row pgsql.Row) (*domain.CarStatistics, error)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func scanRow(row pgsql.Row) (*domain.Car, error) {
	f := &domain.Car{}

	err := row.Scan(&f.ID, &f.LicensePlate, &f.Make, &f.Model)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return f, nil
}

func scanRow2(row pgsql.Row) (*domain.CarTrip, error) {
	f := &domain.CarTrip{}

	err := row.Scan(&f.FirstName, &f.LastName, &f.Duration, &f.Price)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return f, nil
}

func scanRow3(row pgsql.Row) (*domain.CarStatistics, error) {
	f := &domain.CarStatistics{}

	err := row.Scan(&f.CarName)
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

func scanRows2(rows *sql.Rows, scan scanFunc2) ([]*domain.CarTrip, error) {
	var fs []*domain.CarTrip

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

func scanRows3(rows *sql.Rows, scan scanFunc3) ([]*domain.CarStatistics, error) {
	var fs []*domain.CarStatistics

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

	err := p.conn.QueryRowContext(ctx, getSingleSQL, carId).Scan(&c.ID, &c.LicensePlate, &c.Make, &c.Model, &c.Color, &c.Latitude, &c.Longitude, &c.MinutePrice, &c.HourPrice, &c.DayPrice, &c.KilometerPrice, &c.AirConditioning, &c.USB, &c.Bluetooth, &c.Navigation, &c.ChildSeat, &c.Gearbox, &c.Fuel)
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
		return nil, err
	}

	return c, nil
}

func (p PgRepo) UpdateCar(ctx context.Context, id int, license_plate string, car_make string, car_model string, car_color string, minute_price float64, hour_price float64, day_price float64, kilometer_price float64, ac bool, usb bool, bluetooth bool, navigation bool, child_seat bool, gearbox int, fuel int) (*domain.Car, error) {
	c := &domain.Car{}
	_, err := p.conn.ExecContext(ctx, updateCarSQL, license_plate, car_make, car_model, car_color, minute_price, hour_price, day_price, kilometer_price, ac, usb, bluetooth, navigation, child_seat, gearbox, fuel, id)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *PgRepo) CarTrips(ctx context.Context, carId int) ([]*domain.CarTrip, error) {
	rows, err := p.conn.QueryContext(ctx, carTripsSQL, carId)
	if err != nil {
		return nil, err
	}

	return scanRows2(rows, scanRow2)
}

func (p *PgRepo) Statistics(ctx context.Context) ([]*domain.CarStatistics, error) {
	rows, err := p.conn.QueryContext(ctx, statisticsSQL)
	if err != nil {
		return nil, err
	}

	return scanRows3(rows, scanRow3)
}
