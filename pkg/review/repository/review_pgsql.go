package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	insertSQL               = "INSERT INTO įvertinimai (fk_uzklausa, žvaigždutės, komentaras, data) VALUES ($1, $2, $3, $4) RETURNING id"
	getByTicketSQL          = "SELECT id, fk_uzklausa, žvaigždutės, komentaras, data FROM įvertinimai WHERE fk_uzklausa = $1 ORDER BY id ASC LIMIT 1"
	getByTicketForUpdateSQL = "SELECT į.id, į.fk_uzklausa, į.žvaigždutės, į.komentaras, į.data FROM užklausos u INNER JOIN įvertinimai į ON (į.fk_uzklausa = u.id) WHERE u.id = $1 ORDER BY į.id ASC LIMIT 1 FOR UPDATE"
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

func (p *PgRepo) insert(ctx context.Context, q pgsql.Querier, rs *domain.Review) error {
	return q.QueryRowContext(ctx, insertSQL, rs.TicketID, rs.Stars, rs.Comment, rs.Time).Scan(&rs.ID)
}

func (p *PgRepo) Insert(ctx context.Context, rs *domain.Review) error {
	return p.insert(ctx, p.conn, rs)
}

func (p *PgRepo) InsertTx(ctx context.Context, tx repository.Transaction, rs *domain.Review) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.insert(ctx, sqlTx, rs)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (p *PgRepo) getByTicket(ctx context.Context, q pgsql.Querier, ticketID int, forUpdate bool) (*domain.Review, error) {
	var query string
	r := &domain.Review{}

	if forUpdate {
		query = getByTicketForUpdateSQL
	} else {
		query = getByTicketSQL
	}

	err := q.QueryRowContext(ctx, query, ticketID).Scan(&r.ID, &r.TicketID, &r.Stars, &r.Comment, &r.Time)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return r, nil
}

func (p *PgRepo) GetByTicket(ctx context.Context, ticketID int) (*domain.Review, error) {
	return p.getByTicket(ctx, p.conn, ticketID, false)
}

func (p *PgRepo) GetByTicketTx(ctx context.Context, tx repository.Transaction, ticketID int) (*domain.Review, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	r, err := p.getByTicket(ctx, sqlTx, ticketID, true)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return nil, err
	}

	return r, nil
}
