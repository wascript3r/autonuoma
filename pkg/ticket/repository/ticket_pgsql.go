package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	insertSQL                     = "INSERT INTO tickets (user_id) VALUES ($1) RETURNING id"
	getCurrTicketIDSQL            = "SELECT id FROM tickets WHERE user_id = $1 ORDER BY id DESC LIMIT 1"
	getCurrTicketIDForUpdateSQL   = getCurrTicketIDSQL + " FOR UPDATE"
	isCurrTicketEndedSQL          = "SELECT CASE WHEN ended IS NULL THEN false ELSE true END AS ended_b FROM tickets WHERE user_id = $1 ORDER BY id DESC LIMIT 1"
	isCurrTicketEndedForUpdateSQL = isCurrTicketEndedSQL + " FOR UPDATE"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) NewTx(ctx context.Context) (repository.Transaction, error) {
	return p.conn.BeginTx(ctx, nil)
}

func (p *PgRepo) insert(ctx context.Context, q pgsql.Querier, ts *domain.Ticket) error {
	return q.QueryRowContext(ctx, insertSQL, ts.UserID).Scan(&ts.ID)
}

func (p *PgRepo) Insert(ctx context.Context, ts *domain.Ticket) error {
	return p.insert(ctx, p.conn, ts)
}

func (p *PgRepo) InsertTx(ctx context.Context, tx repository.Transaction, ts *domain.Ticket) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.insert(ctx, sqlTx, ts)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (p *PgRepo) getCurrTicketID(ctx context.Context, q pgsql.Querier, userID int, forUpdate bool) (int, error) {
	var (
		ticketID int
		query    string
	)

	if forUpdate {
		query = getCurrTicketIDForUpdateSQL
	} else {
		query = getCurrTicketIDSQL
	}

	err := q.QueryRowContext(ctx, query, userID).Scan(&ticketID)
	if err != nil {
		return 0, pgsql.ParseSQLError(err)
	}

	return ticketID, nil
}

func (p *PgRepo) GetCurrTicketID(ctx context.Context, userID int) (int, error) {
	return p.getCurrTicketID(ctx, p.conn, userID, false)
}

func (p *PgRepo) GetCurrTicketIDTx(ctx context.Context, tx repository.Transaction, userID int) (int, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return 0, repository.ErrTxMismatch
	}

	ticketID, err := p.getCurrTicketID(ctx, sqlTx, userID, true)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return 0, err
	}

	return ticketID, nil
}

func (p *PgRepo) isCurrTicketEnded(ctx context.Context, q pgsql.Querier, userID int, forUpdate bool) (bool, error) {
	var (
		ended bool
		query string
	)

	if forUpdate {
		query = isCurrTicketEndedForUpdateSQL
	} else {
		query = isCurrTicketEndedSQL
	}

	err := q.QueryRowContext(ctx, query, userID).Scan(&ended)
	if err != nil {
		return false, pgsql.ParseSQLError(err)
	}

	return ended, nil
}

func (p *PgRepo) IsCurrTicketEnded(ctx context.Context, userID int) (bool, error) {
	return p.isCurrTicketEnded(ctx, p.conn, userID, false)
}

func (p *PgRepo) IsCurrTicketEndedTx(ctx context.Context, tx repository.Transaction, userID int) (bool, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return false, repository.ErrTxMismatch
	}

	ended, err := p.isCurrTicketEnded(ctx, sqlTx, userID, true)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return false, err
	}

	return ended, nil
}
