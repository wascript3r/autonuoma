package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	insertSQL = "INSERT INTO sessions (id, user_id, expiration) VALUES ($1, $2, $3)"
	getSQL    = "SELECT s.*, COALESCE(r.is_admin, false) AS is_admin FROM sessions s LEFT JOIN roles r ON r.user_id = s.user_id WHERE s.id = $1"
	deleteSQL = "DELETE FROM sessions WHERE id = $1"
)

type PgRepo struct {
	conn *sql.DB
}

func NewPgRepo(c *sql.DB) *PgRepo {
	return &PgRepo{c}
}

func (p *PgRepo) Insert(ctx context.Context, ss *domain.Session) error {
	_, err := p.conn.ExecContext(ctx, insertSQL, ss.ID, ss.UserID, ss.Expiration)
	return err
}

func (p PgRepo) Get(ctx context.Context, id string) (*domain.Session, error) {
	s := &domain.Session{}

	err := p.conn.QueryRowContext(ctx, getSQL, id).Scan(&s.ID, &s.UserID, &s.Expiration, &s.Role)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return s, nil
}

func (p PgRepo) Delete(ctx context.Context, id string) error {
	_, err := p.conn.ExecContext(ctx, deleteSQL, id)
	return err
}
