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
	insertSQL        = "INSERT INTO užklausos (fk_klientas, fk_klientų_aptarnavimo_specialistas, sukurta, užbaigta) VALUES ($1, $2, $3, $4) RETURNING id"
	setAgentSQL      = "UPDATE užklausos SET fk_klientų_aptarnavimo_specialistas = $2 WHERE id = $1"
	setEndedSQL      = "UPDATE užklausos SET užbaigta = $2 WHERE id = $1"
	setAgentEndedSQL = "UPDATE užklausos SET fk_klientų_aptarnavimo_specialistas = $2, užbaigta = $3 WHERE id = $1"

	getLastActiveTicketIDSQL          = "SELECT id FROM užklausos WHERE fk_klientas = $1 AND užbaigta IS NULL ORDER BY id DESC LIMIT 1"
	getLastActiveTicketIDForUpdateSQL = getLastActiveTicketIDSQL + " FOR UPDATE"

	getTicketMetaSQL          = "SELECT fk_klientas, fk_klientų_aptarnavimo_specialistas, užbaigta FROM užklausos WHERE id = $1"
	getTicketMetaForUpdateSQL = getTicketMetaSQL + " FOR UPDATE"

	getTicketsSQL     = "SELECT u.id, u.fk_klientų_aptarnavimo_specialistas, u.užbaigta, v.id, v.vardas, v.pavardė, ž.tekstas, ž.išsiųsta FROM užklausos u INNER JOIN vartotojai v ON (v.id = u.fk_klientas) INNER JOIN (SELECT fk_uzklausa, tekstas, išsiųsta FROM žinutės WHERE id IN (SELECT MIN(id) FROM žinutės GROUP BY fk_uzklausa)) ž ON (ž.fk_uzklausa = u.id) ORDER BY u.id ASC"
	getUserTicketsSQL = "SELECT u.id, u.fk_klientų_aptarnavimo_specialistas, u.užbaigta, v.id, v.vardas, v.pavardė, ž.tekstas, ž.išsiųsta FROM užklausos u INNER JOIN vartotojai v ON (v.id = u.fk_klientas) INNER JOIN (SELECT fk_uzklausa, tekstas, išsiųsta FROM žinutės WHERE id IN (SELECT MIN(id) FROM žinutės GROUP BY fk_uzklausa)) ž ON (ž.fk_uzklausa = u.id) WHERE u.fk_klientas = $1 ORDER BY u.id ASC"
)

type scanFunc func(row pgsql.Row) (*domain.TicketFull, error)

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
	return q.QueryRowContext(ctx, insertSQL, ts.ClientID, ts.AgentID, ts.Created, ts.Ended).Scan(&ts.ID)
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

func (p *PgRepo) setAgent(ctx context.Context, q pgsql.Querier, id int, agentID int) error {
	_, err := q.ExecContext(ctx, setAgentSQL, id, agentID)
	return err
}

func (p *PgRepo) SetAgent(ctx context.Context, id int, agentID int) error {
	return p.setAgent(ctx, p.conn, id, agentID)
}

func (p *PgRepo) SetAgentTx(ctx context.Context, tx repository.Transaction, id int, agentID int) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.setAgent(ctx, sqlTx, id, agentID)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (p *PgRepo) setEnded(ctx context.Context, q pgsql.Querier, id int, ended time.Time) error {
	_, err := q.ExecContext(ctx, setEndedSQL, id, ended)
	return err
}

func (p *PgRepo) SetEnded(ctx context.Context, id int, ended time.Time) error {
	return p.setEnded(ctx, p.conn, id, ended)
}

func (p *PgRepo) SetEndedTx(ctx context.Context, tx repository.Transaction, id int, ended time.Time) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.setEnded(ctx, sqlTx, id, ended)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (p *PgRepo) setAgentEnded(ctx context.Context, q pgsql.Querier, id int, agentID int, ended time.Time) error {
	_, err := q.ExecContext(ctx, setAgentEndedSQL, id, agentID, ended)
	return err
}

func (p *PgRepo) SetAgentEnded(ctx context.Context, id int, agentID int, ended time.Time) error {
	return p.setAgentEnded(ctx, p.conn, id, agentID, ended)
}

func (p *PgRepo) SetAgentEndedTx(ctx context.Context, tx repository.Transaction, id int, agentID int, ended time.Time) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.setAgentEnded(ctx, sqlTx, id, agentID, ended)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (p *PgRepo) getLastActiveTicketID(ctx context.Context, q pgsql.Querier, clientID int, forUpdate bool) (int, error) {
	var (
		ticketID int
		query    string
	)

	if forUpdate {
		query = getLastActiveTicketIDForUpdateSQL
	} else {
		query = getLastActiveTicketIDSQL
	}

	err := q.QueryRowContext(ctx, query, clientID).Scan(&ticketID)
	if err != nil {
		return 0, pgsql.ParseSQLError(err)
	}

	return ticketID, nil
}

func (p *PgRepo) GetLastActiveTicketID(ctx context.Context, clientID int) (int, error) {
	return p.getLastActiveTicketID(ctx, p.conn, clientID, false)
}

func (p *PgRepo) GetLastActiveTicketIDTx(ctx context.Context, tx repository.Transaction, clientID int) (int, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return 0, repository.ErrTxMismatch
	}

	ticketID, err := p.getLastActiveTicketID(ctx, sqlTx, clientID, true)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return 0, err
	}

	return ticketID, nil
}

func getStatus(agentID *int, ended *time.Time) domain.TicketStatus {
	if ended == nil {
		if agentID == nil {
			return domain.CreatedTicketStatus
		}
		return domain.AcceptedTicketStatus
	}
	return domain.EndedTicketStatus
}

func (p *PgRepo) getTicketMeta(ctx context.Context, q pgsql.Querier, id int, forUpdate bool) (*domain.TicketMeta, error) {
	var query string
	m := &domain.TicketMeta{}

	if forUpdate {
		query = getTicketMetaForUpdateSQL
	} else {
		query = getTicketMetaSQL
	}

	err := q.QueryRowContext(ctx, query, id).Scan(&m.ClientID, &m.AgentID, &m.Ended)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	m.Status = getStatus(m.AgentID, m.Ended)
	return m, nil
}

func (p *PgRepo) GetTicketMeta(ctx context.Context, id int) (*domain.TicketMeta, error) {
	return p.getTicketMeta(ctx, p.conn, id, false)
}

func (p *PgRepo) GetTicketMetaTx(ctx context.Context, tx repository.Transaction, id int) (*domain.TicketMeta, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	meta, err := p.getTicketMeta(ctx, sqlTx, id, true)
	if err != nil {
		sqlTx.Rollback()
		return nil, err
	}

	return meta, nil
}

func scanRow(row pgsql.Row) (*domain.TicketFull, error) {
	var (
		agentID *int
		ended   *time.Time
	)

	t := &domain.TicketFull{
		ID:           0,
		Status:       0,
		ClientMeta:   &domain.UserMeta{},
		FirstMessage: "",
		Time:         time.Time{},
	}

	err := row.Scan(
		&t.ID,
		&agentID,
		&ended,

		&t.ClientMeta.ID,
		&t.ClientMeta.FirstName,
		&t.ClientMeta.LastName,

		&t.FirstMessage,
		&t.Time,
	)
	if err != nil {
		return nil, pgsql.ParseSQLError(err)
	}

	t.Status = getStatus(agentID, ended)
	return t, nil
}

func scanRows(rows *sql.Rows, scan scanFunc) ([]*domain.TicketFull, error) {
	var ts []*domain.TicketFull

	for rows.Next() {
		m, err := scan(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		ts = append(ts, m)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return ts, nil
}

func (p *PgRepo) getTickets(ctx context.Context, q pgsql.Querier, userID *int) ([]*domain.TicketFull, error) {
	var (
		query string
		args  []interface{}
	)
	if userID == nil {
		query = getTicketsSQL
	} else {
		query = getUserTicketsSQL
		args = append(args, *userID)
	}

	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return scanRows(rows, scanRow)
}

func (p *PgRepo) GetTickets(ctx context.Context) ([]*domain.TicketFull, error) {
	return p.getTickets(ctx, p.conn, nil)
}

func (p *PgRepo) GetTicketsTx(ctx context.Context, tx repository.Transaction) ([]*domain.TicketFull, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	ms, err := p.getTickets(ctx, sqlTx, nil)
	if err != nil {
		sqlTx.Rollback()
		return nil, err
	}

	return ms, nil
}

func (p *PgRepo) GetUserTickets(ctx context.Context, userID int) ([]*domain.TicketFull, error) {
	return p.getTickets(ctx, p.conn, &userID)
}

func (p *PgRepo) GetUserTicketsTx(ctx context.Context, tx repository.Transaction, userID int) ([]*domain.TicketFull, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	ms, err := p.getTickets(ctx, sqlTx, &userID)
	if err != nil {
		sqlTx.Rollback()
		return nil, err
	}

	return ms, nil
}
