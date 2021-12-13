package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
	"github.com/wascript3r/autonuoma/pkg/user"
)

const (
	insertIfNotExistsSQL = "INSERT INTO vartotojai (vardas, pavardė, el_paštas, gimimo_data, slaptažodis, balansas, asmens_kodas, rolė) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	emailExistsSQL       = "SELECT EXISTS(SELECT 1 FROM vartotojai WHERE el_paštas = $1)"
	getCredentialsSQL    = "SELECT id, rolė, slaptažodis FROM vartotojai WHERE el_paštas = $1"
	deductBalanceSQL     = "UPDATE vartotojai SET balansas = balansas - $2 WHERE id = $1"
	addBalanceSQL        = "UPDATE vartotojai SET balansas = balansas + $2 WHERE id = $1"
	getDataSQL           = "SELECT vardas, pavardė, el_paštas, gimimo_data, balansas FROM vartotojai WHERE id = $1"
	getLicenseStatusSQL  = "SELECT b.name, p.galiojimo_pabaiga FROM vairuotojo_pažymėjimai p, vairuotojo_pažymėjimo_būsenos b WHERE p.fk_vartotojas = $1 AND b.id = p.būsena AND p.būsena = $2"
	updateEmailSQL       = "UPDATE vartotojai SET el_paštas = $2 WHERE id = $1"
	updatePasswordSQL    = "UPDATE vartotojai SET slaptažodis = $2 WHERE id = $1"
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

		us.FirstName,
		us.LastName,
		us.Email,
		us.BirthDate,
		us.Password,
		us.Balance,
		us.PIN,
		us.RoleID,
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

func (p PgRepo) GetCredentials(ctx context.Context, email string) (*domain.UserCredentials, error) {
	c := &domain.UserCredentials{}

	err := p.conn.QueryRowContext(ctx, getCredentialsSQL, email).Scan(&c.ID, &c.RoleID, &c.Password)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return c, nil
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

func (p *PgRepo) GetLicenseStatus(ctx context.Context, uid int) (string, error) {
	var licenseStatus string
	var licenseEndDate time.Time

	if err := p.conn.QueryRowContext(ctx, getLicenseStatusSQL, uid, domain.ConfirmedLicenseStatus).Scan(&licenseStatus, &licenseEndDate); err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if len(licenseStatus) > 0 {
		if licenseEndDate.Before(time.Now()) {
			return "pasibaigęs galiojimas", nil
		}

		return strings.TrimSpace(licenseStatus), nil
	}

	if err := p.conn.QueryRowContext(ctx, getLicenseStatusSQL, uid, domain.SubmittedLicenseStatus).Scan(&licenseStatus, &licenseEndDate); err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if len(licenseStatus) > 0 {
		return strings.TrimSpace(licenseStatus), nil
	}

	if err := p.conn.QueryRowContext(ctx, getLicenseStatusSQL, uid, domain.RejectedLicenseStatus).Scan(&licenseStatus, &licenseEndDate); err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if len(licenseStatus) > 0 {
		return strings.TrimSpace(licenseStatus), nil
	}

	return "nepateiktas", nil
}

func (p *PgRepo) GetData(ctx context.Context, uid int) (*user.UserProfile, error) {
	u := &user.UserProfile{}
	u.ID = uid

	err := p.conn.QueryRowContext(ctx, getDataSQL, uid).Scan(&u.FirstName, &u.LastName, &u.Email, &u.Birthdate, &u.Balance)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return u, nil
}

func (p *PgRepo) UpdateEmail(ctx context.Context, uid int, email string) error {
	if err := p.conn.QueryRowContext(ctx, updateEmailSQL, uid, email).Err(); err != nil {
		return pgsql.ParseSQLError(err)
	}
	return nil
}

func (p *PgRepo) UpdatePassword(ctx context.Context, uid int, hash string) error {
	if err := p.conn.QueryRowContext(ctx, updatePasswordSQL, uid, hash).Err(); err != nil {
		return pgsql.ParseSQLError(err)
	}
	return nil
}
