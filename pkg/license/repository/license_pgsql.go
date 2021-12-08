package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	getStatusSQL          = "SELECT būsena FROM vairuotojo_pažymėjimai WHERE id = $1"
	getStatusForUpdateSQL = getStatusSQL + " FOR UPDATE"

	setStatusSQL = "UPDATE vairuotojo_pažymėjimai SET būsena = $2 WHERE id = $1"
	getAllSQL    = "SELECT vp.id, vp.nr, v.id, v.vardas, v.pavardė, vp.galiojimo_pabaiga, vp.būsena FROM vairuotojo_pažymėjimai vp INNER JOIN vartotojai v ON (v.id = vp.fk_vartotojas) WHERE vp.būsena = $1 ORDER BY vp.id ASC"
	getPhotosSQL = "SELECT id, fk_vairuotojo_pazymejimas, nuoroda FROM vairuotojo_pažymėjimo_nuotraukos WHERE fk_vairuotojo_pazymejimas = $1 ORDER BY id ASC"
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

func scanLicense(row pgsql.Row) (*domain.LicenseFull, error) {
	l := &domain.LicenseFull{
		ID:         0,
		Number:     "",
		ClientMeta: &domain.UserMeta{},
		Expiration: time.Time{},
		StatusID:   0,
	}

	err := row.Scan(
		&l.ID,
		&l.Number,

		&l.ClientMeta.ID,
		&l.ClientMeta.FirstName,
		&l.ClientMeta.LastName,

		&l.Expiration,
		&l.StatusID,
	)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return l, nil
}

func scanLicenses(rows *sql.Rows, scan func(row pgsql.Row) (*domain.LicenseFull, error)) ([]*domain.LicenseFull, error) {
	var ls []*domain.LicenseFull

	for rows.Next() {
		l, err := scan(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		ls = append(ls, l)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return ls, nil
}

func (p *PgRepo) getAll(ctx context.Context, q pgsql.Querier) ([]*domain.LicenseFull, error) {
	rows, err := q.QueryContext(ctx, getAllSQL, domain.SubmittedLicenseStatus)
	if err != nil {
		return nil, err
	}

	return scanLicenses(rows, scanLicense)
}

func (p *PgRepo) GetAllUnconfirmed(ctx context.Context) ([]*domain.LicenseFull, error) {
	return p.getAll(ctx, p.conn)
}

func (p *PgRepo) GetAllUnconfirmedTx(ctx context.Context, tx repository.Transaction) ([]*domain.LicenseFull, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	ls, err := p.getAll(ctx, sqlTx)
	if err != nil {
		sqlTx.Rollback()
		return nil, err
	}

	return ls, nil
}

func scanPhoto(row pgsql.Row) (*domain.LicensePhoto, error) {
	p := &domain.LicensePhoto{}

	err := row.Scan(&p.ID, &p.LicenseID, &p.URL)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return p, nil
}

func scanPhotos(rows *sql.Rows, scan func(row pgsql.Row) (*domain.LicensePhoto, error)) ([]*domain.LicensePhoto, error) {
	var ps []*domain.LicensePhoto

	for rows.Next() {
		p, err := scan(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		ps = append(ps, p)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return ps, nil
}

func (p *PgRepo) getPhotos(ctx context.Context, q pgsql.Querier, licenseID int) ([]*domain.LicensePhoto, error) {
	rows, err := q.QueryContext(ctx, getPhotosSQL, licenseID)
	if err != nil {
		return nil, err
	}

	return scanPhotos(rows, scanPhoto)
}

func (p *PgRepo) GetPhotos(ctx context.Context, licenseID int) ([]*domain.LicensePhoto, error) {
	return p.getPhotos(ctx, p.conn, licenseID)
}

func (p *PgRepo) GetPhotosTx(ctx context.Context, tx repository.Transaction, licenseID int) ([]*domain.LicensePhoto, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	ps, err := p.getPhotos(ctx, sqlTx, licenseID)
	if err != nil {
		sqlTx.Rollback()
		return nil, err
	}

	return ps, nil
}
