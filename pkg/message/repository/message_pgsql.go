package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	insertSQL = "INSERT INTO žinutės (fk_uzklausa, fk_vartotojas, tekstas, išsiųsta) VALUES ($1, $2, $3, $4) RETURNING id"
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

func (p *PgRepo) insert(ctx context.Context, q pgsql.Querier, ms *domain.Message) error {
	return q.QueryRowContext(ctx, insertSQL, ms.TicketID, ms.UserID, ms.Content, ms.Time).Scan(&ms.ID)
}

func (p *PgRepo) Insert(ctx context.Context, ms *domain.Message) error {
	return p.insert(ctx, p.conn, ms)
}

func (p *PgRepo) InsertTx(ctx context.Context, tx repository.Transaction, ms *domain.Message) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.insert(ctx, sqlTx, ms)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}
