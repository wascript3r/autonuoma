package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
	"time"
)

const (
	startTripSQL   = "INSERT INTO kelionės (pradžios_laikas, pabaigos_taško_ilguma, pabaigos_taško_platuma, fk_rezervacija, kaina) VALUES ($1, $2, $3, $4, 0) RETURNING id, pradžios_laikas"
	endTripSQL     = "UPDATE kelionės SET pabaigos_laikas = $2, kaina = $3 WHERE id = $1"
	getTripByIdSQL = "SELECT pradžios_laikas, pabaigos_taško_platuma, pabaigos_taško_ilguma, id, fk_rezervacija FROM kelionės WHERE fk_rezervacija = $1"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) Start(ctx context.Context, endLng string, endLat string, reservationID int) (int, time.Time, error) {
	fmt.Println("trip start")
	var tripID int
	var createdAt time.Time
	err := p.conn.QueryRowContext(ctx, startTripSQL, time.Now(), endLng, endLat, reservationID).Scan(&tripID, &createdAt)

	fmt.Println(err)
	if err != nil {
		return 0, time.Now(), err
	}
	return tripID, createdAt, nil
}

func (p *PgRepo) End(ctx context.Context, tripID int, price float32) error {
	if err := p.conn.QueryRowContext(ctx, endTripSQL, tripID, time.Now(), price).Err(); err != nil {
		return pgsql.ParseSQLError(err)
	}
	return nil
}

func (p *PgRepo) GetByReservationId(ctx context.Context, reservationID int) (*domain.Trip, error) {
	t := &domain.Trip{}
	t.ReservationID = reservationID

	err := p.conn.QueryRowContext(ctx, getTripByIdSQL, reservationID).Scan(&t.Begin, &t.EndLat, &t.EndLng, &t.ID, &t.ReservationID)
	fmt.Println(err)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return t, nil
}
