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
	insertSQL            = "INSERT INTO žinutės (fk_uzklausa, fk_vartotojas, tekstas, išsiųsta) VALUES ($1, $2, $3, $4) RETURNING id"
	getTicketMessagesSQL = "SELECT v.id, v.vardas, v.pavardė, ž.tekstas, ž.išsiųsta FROM užklausos u INNER JOIN žinutės ž ON (ž.fk_uzklausa = u.id) INNER JOIN vartotojai v ON (v.id = ž.fk_vartotojas) WHERE u.id = $1 ORDER BY ž.id ASC"
)

type scanFunc func(row pgsql.Row) (*domain.MessageFull, error)

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

func scanRow(row pgsql.Row) (*domain.MessageFull, error) {
	m := &domain.MessageFull{
		UserMeta: &domain.UserMeta{},
		Content:  "",
		Time:     time.Time{},
	}

	err := row.Scan(
		&m.UserMeta.ID,
		&m.UserMeta.FirstName,
		&m.UserMeta.LastName,

		&m.Content,
		&m.Time,
	)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	return m, nil
}

func scanRows(rows *sql.Rows, scan scanFunc) ([]*domain.MessageFull, error) {
	var ms []*domain.MessageFull

	for rows.Next() {
		m, err := scan(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		ms = append(ms, m)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return ms, nil
}

func (p *PgRepo) getTicketMessages(ctx context.Context, q pgsql.Querier, ticketID int) ([]*domain.MessageFull, error) {
	rows, err := q.QueryContext(ctx, getTicketMessagesSQL, ticketID)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, scanRow)
}

func (p *PgRepo) GetTicketMessages(ctx context.Context, ticketID int) ([]*domain.MessageFull, error) {
	return p.getTicketMessages(ctx, p.conn, ticketID)
}

func (p *PgRepo) GetTicketMessagesTx(ctx context.Context, tx repository.Transaction, ticketID int) ([]*domain.MessageFull, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	ms, err := p.getTicketMessages(ctx, sqlTx, ticketID)
	if err != nil {
		sqlTx.Rollback()
		return nil, err
	}

	return ms, nil
}
