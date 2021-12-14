package repository

import (
	"context"
	"database/sql"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
	"time"
)

const (
	createReservationSQL     = "INSERT INTO rezervacijos (sukurta, fk_automobilis, fk_vairuotojas) VALUES ($1, $2, $3) RETURNING id"
	cancelReservationSQL     = "UPDATE rezervacijos SET atšaukta = $2 WHERE id = $1"
	getCurrentReservationSQl = "SELECT sukurta, atšaukta, pradzios_adresas, pabaigos_adresas, id, fk_automobilis, fk_vartotojas FROM rezervacijos WHERE atšaukta = NULL AND fk_vartotojas = $1"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) Create(ctx context.Context, carID int, userID int) (int, error) {
	var reservationID int
	err := p.conn.QueryRowContext(ctx, createReservationSQL, time.Now(), carID, userID).Scan(&reservationID)
	if err != nil {
		return 0, err
	}
	return reservationID, nil
}

func (p *PgRepo) Cancel(ctx context.Context, reservationID int) error {
	if err := p.conn.QueryRowContext(ctx, cancelReservationSQL, reservationID, time.Now()).Err(); err != nil {
		return pgsql.ParseSQLError(err)
	}
	return nil
}

func (p *PgRepo) GetCurrent(ctx context.Context, userID int) (*domain.Reservation, error) {
	r := &domain.Reservation{}

	err := p.conn.QueryRowContext(ctx, getCurrentReservationSQl, userID).Scan(&r.CreatedAt, &r.CanceledAt, &r.StartAddress, &r.EndAddress, &r.ID, &r.CarID, &r.UserID)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return r, nil
}
