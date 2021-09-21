package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	insertIfNotExistsSQL = "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	emailExistsSQL       = "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	getIDAndPasswordSQL  = "SELECT id, password FROM users WHERE email = $1"
	deductBalanceSQL     = "UPDATE users SET balance = balance - $2 WHERE id = $1"
	addBalanceSQL        = "UPDATE users SET balance = balance + $2 WHERE id = $1"
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

func (p *PgRepo) InsertIfNotExists(ctx context.Context, us *domain.User) error {
	err := p.conn.QueryRowContext(
		ctx,
		insertIfNotExistsSQL,

		us.Email,
		us.Password,
	).Scan(&us.ID)

	return pgsql.ParsePgError(err)
}

func (p PgRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := p.conn.QueryRowContext(ctx, emailExistsSQL, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p PgRepo) GetIDAndPassword(ctx context.Context, email string) (int, string, error) {
	var (
		id       int
		password string
	)

	err := p.conn.QueryRowContext(ctx, getIDAndPasswordSQL, email).Scan(&id, &password)
	if err != nil {
		return 0, "", pgsql.ParseSQLError(err)
	}

	return id, password, nil
}

func (p *PgRepo) deductBalance(ctx context.Context, q pgsql.Querier, id int, value int64) error {
	res, err := q.ExecContext(ctx, deductBalanceSQL, id, value)
	if err != nil {
		return pgsql.ParsePgError(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (p *PgRepo) DeductBalance(ctx context.Context, id int, value int64) error {
	return p.deductBalance(ctx, p.conn, id, value)
}

func (p *PgRepo) DeductBalanceTx(ctx context.Context, tx repository.Transaction, id int, value int64) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.deductBalance(ctx, sqlTx, id, value)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (p *PgRepo) addBalance(ctx context.Context, q pgsql.Querier, id int, value int64) error {
	res, err := q.ExecContext(ctx, addBalanceSQL, id, value)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (p *PgRepo) AddBalance(ctx context.Context, id int, value int64) error {
	return p.addBalance(ctx, p.conn, id, value)
}

func (p *PgRepo) AddBalanceTx(ctx context.Context, tx repository.Transaction, id int, value int64) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.addBalance(ctx, sqlTx, id, value)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}
