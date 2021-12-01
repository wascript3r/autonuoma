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
	insertSQL   = "INSERT INTO užklausos (fk_klientas, fk_klientų_aptarnavimo_specialistas, sukurta, užbaigta) VALUES ($1, $2, $3, $4) RETURNING id"
	setEndedSQL = "UPDATE užklausos SET užbaigta = $2 WHERE id = $1"
	setAgentSQL = "UPDATE užklausos SET fk_klientų_aptarnavimo_specialistas = $2 WHERE id = $1"

	getLastActiveTicketIDSQL          = "SELECT id FROM užklausos WHERE fk_klientas = $1 AND užbaigta IS NULL ORDER BY id DESC LIMIT 1"
	getLastActiveTicketIDForUpdateSQL = getLastActiveTicketIDSQL + " FOR UPDATE"

	getTicketStatusSQL          = "SELECT užbaigta, fk_klientų_aptarnavimo_specialistas FROM užklausos WHERE id = $1"
	getTicketStatusForUpdateSQL = getTicketStatusSQL + " FOR UPDATE"
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

func (p *PgRepo) getTicketStatus(ctx context.Context, q pgsql.Querier, id int, forUpdate bool) (domain.TicketStatus, error) {
	var (
		query   string
		ended   *time.Time
		agentID *int
	)

	if forUpdate {
		query = getTicketStatusForUpdateSQL
	} else {
		query = getTicketStatusSQL
	}

	err := q.QueryRowContext(ctx, query, id).Scan(&ended, &agentID)
	if err != nil {
		return 0, pgsql.ParseSQLError(err)
	}

	if ended == nil {
		if agentID == nil {
			return domain.CreatedTicketStatus, nil
		} else {
			return domain.AcceptedTicketStatus, nil
		}
	}

	return domain.EndedTicketStatus, nil
}

func (p *PgRepo) GetTicketStatus(ctx context.Context, id int) (domain.TicketStatus, error) {
	return p.getTicketStatus(ctx, p.conn, id, false)
}

func (p *PgRepo) GetTicketStatusTx(ctx context.Context, tx repository.Transaction, id int) (domain.TicketStatus, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return 0, repository.ErrTxMismatch
	}

	status, err := p.getTicketStatus(ctx, sqlTx, id, true)
	if err != nil {
		sqlTx.Rollback()
		return 0, err
	}

	return status, nil
}
