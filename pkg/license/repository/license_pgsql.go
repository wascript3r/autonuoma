package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	getStatusSQL          = "SELECT būsena FROM vairuotojo_pažymėjimai WHERE id = $1"
	getStatusForUpdateSQL = getStatusSQL + " FOR UPDATE"

	setStatusSQL = "UPDATE vairuotojo_pažymėjimai SET būsena = $2 WHERE id = $1"
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

func (p *PgRepo) getStatus(ctx context.Context, q pgsql.Querier, id int, forUpdate bool) (domain.LicenseStatus, error) {
	var (
		query  string
		status domain.LicenseStatus
	)

	if forUpdate {
		query = getStatusForUpdateSQL
	} else {
		query = getStatusSQL
	}

	err := q.QueryRowContext(ctx, query, id).Scan(&status)
	if err != nil {
		return 0, pgsql.ParseSQLError(err)
	}

	return status, nil
}

func (p *PgRepo) GetStatus(ctx context.Context, id int) (domain.LicenseStatus, error) {
	return p.getStatus(ctx, p.conn, id, false)
}

func (p *PgRepo) GetStatusTx(ctx context.Context, tx repository.Transaction, id int) (domain.LicenseStatus, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return 0, repository.ErrTxMismatch
	}

	status, err := p.getStatus(ctx, sqlTx, id, true)
	if err != nil {
		sqlTx.Rollback()
		return 0, err
	}

	return status, nil
}

func (p *PgRepo) setStatus(ctx context.Context, q pgsql.Querier, id int, status domain.LicenseStatus) error {
	_, err := q.ExecContext(ctx, setStatusSQL, id, status)
	return err
}

func (p *PgRepo) SetStatus(ctx context.Context, id int, status domain.LicenseStatus) error {
	return p.setStatus(ctx, p.conn, id, status)
}

func (p *PgRepo) SetStatusTx(ctx context.Context, tx repository.Transaction, id int, status domain.LicenseStatus) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.setStatus(ctx, sqlTx, id, status)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}
