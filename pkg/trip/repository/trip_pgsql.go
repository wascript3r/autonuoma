package repository

import (
	"context"
	"database/sql"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
	"time"
)

const (
	startTripSQL   = "INSERT INTO kelionės (pradžios_laikas, fk_rezervacija) VALUES ($1, $2) RETURNING id"
	endTripSQL     = "UPDATE kelionės SET pabaigos_laikas = $2, pabaigos_taško_platuma = $3, pabaigos_taško_ilguma = $4, trukmė = $5 WHERE id = $1"
	getTripByIdSQL = "SELECT pradžios_laikas, pabaigos_laikas, pabaigos_taško_platuma, pabaigos_taško_ilguma, trukmė, kaina, id, fk_rezervacija FROM kelionės WHERE fk_rezervacija = $1"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) Start(ctx context.Context, startTime time.Time, reservationID int) (int, error) {
	var tripID int
	err := p.conn.QueryRowContext(ctx, startTripSQL, startTime, reservationID).Scan(&tripID)
	if err != nil {
		return 0, err
	}
	return tripID, nil
}

func (p *PgRepo) End(ctx context.Context, tripID int, endLat string, endLng string) error {
	if err := p.conn.QueryRowContext(ctx, endTripSQL, tripID, time.Now(), endLat, endLng).Err(); err != nil {
		return pgsql.ParseSQLError(err)
	}
	return nil
}

func (p *PgRepo) GetByReservationId(ctx context.Context, reservationID int) (*domain.Trip, error) {
	t := &domain.Trip{}
	t.ReservationID = reservationID

	err := p.conn.QueryRowContext(ctx, getTripByIdSQL, reservationID).Scan(&t.Begin, &t.End, &t.EndLat, &t.EndLng, &t.Duration, &t.Price, &t.ID)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return t, nil
}
