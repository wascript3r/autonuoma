package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	getAllSQL = "SELECT id, kategorija, klausimas, atsakymas FROM dažniausiai_užduodami_klausimai ORDER BY id ASC"
)

type scanFunc = func(row pgsql.Row) (*domain.FAQ, error)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func scanRow(row pgsql.Row) (*domain.FAQ, error) {
	f := &domain.FAQ{}

	err := row.Scan(&f.ID, &f.CategoryID, &f.Question, &f.Answer)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return f, nil
}

func scanRows(rows *sql.Rows, scan scanFunc) ([]*domain.FAQ, error) {
	var fs []*domain.FAQ

	for rows.Next() {
		f, err := scan(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		fs = append(fs, f)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return fs, nil
}

func (p *PgRepo) GetAll(ctx context.Context) ([]*domain.FAQ, error) {
	rows, err := p.conn.QueryContext(ctx, getAllSQL)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, scanRow)
}
