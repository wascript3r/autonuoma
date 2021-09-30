package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
	"github.com/wascript3r/autonuoma/pkg/repository/pgsql"
)

const (
	insertIfNotExistsSQL          = "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	emailExistsSQL                = "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	getIDAndPasswordSQL           = "SELECT id, password FROM users WHERE email = $1"
	deductBalanceSQL              = "UPDATE users SET balance = balance - $2 WHERE id = $1"
	addBalanceSQL                 = "UPDATE users SET balance = balance + $2 WHERE id = $1"
	getCurrTicketIDSQL            = "SELECT current_ticket_id FROM users WHERE id = $1"
	getCurrTicketIDForUpdateSQL   = getCurrTicketIDSQL + " FOR UPDATE"
	isCurrTicketEndedSQL          = "SELECT t.ended FROM users u LEFT JOIN tickets t ON (t.id = u.current_ticket_id) WHERE u.id = $1"
	isCurrTicketEndedForUpdateSQL = isCurrTicketEndedSQL + " FOR UPDATE OF u"
	setCurrTicketSQL              = "UPDATE users SET current_ticket_id = $2 WHERE id = $1"
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

func (p *PgRepo) getCurrTicketID(ctx context.Context, q pgsql.Querier, userID int, forUpdate bool) (int, error) {
	var (
		ticketID *int
		query    string
	)

	if forUpdate {
		query = getCurrTicketIDForUpdateSQL
	} else {
		query = getCurrTicketIDSQL
	}

	err := q.QueryRowContext(ctx, query, userID).Scan(&ticketID)
	if err != nil {
		return 0, err
	}

	if ticketID == nil {
		return 0, domain.ErrNullValue
	}

	return *ticketID, nil
}

func (p *PgRepo) GetCurrTicketID(ctx context.Context, userID int) (int, error) {
	return p.getCurrTicketID(ctx, p.conn, userID, false)
}

func (p *PgRepo) GetCurrTicketIDTx(ctx context.Context, tx repository.Transaction, userID int) (int, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return 0, repository.ErrTxMismatch
	}

	ticketID, err := p.getCurrTicketID(ctx, sqlTx, userID, true)
	if err != nil {
		if err != domain.ErrNullValue {
			sqlTx.Rollback()
		}
		return 0, err
	}

	return ticketID, nil
}

func (p *PgRepo) isCurrTicketEnded(ctx context.Context, q pgsql.Querier, userID int, forUpdate bool) (bool, error) {
	var (
		ended *bool
		query string
	)

	if forUpdate {
		query = isCurrTicketEndedForUpdateSQL
	} else {
		query = isCurrTicketEndedSQL
	}

	err := q.QueryRowContext(ctx, query, userID).Scan(&ended)
	if err != nil {
		return false, err
	}

	if ended == nil {
		return false, domain.ErrNullValue
	}

	return *ended, nil
}

func (p *PgRepo) IsCurrTicketEnded(ctx context.Context, userID int) (bool, error) {
	return p.isCurrTicketEnded(ctx, p.conn, userID, false)
}

func (p *PgRepo) IsCurrTicketEndedTx(ctx context.Context, tx repository.Transaction, userID int) (bool, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return false, repository.ErrTxMismatch
	}

	ended, err := p.isCurrTicketEnded(ctx, sqlTx, userID, true)
	if err != nil {
		if err != domain.ErrNullValue {
			sqlTx.Rollback()
		}
		return false, err
	}

	return ended, nil
}

func (p *PgRepo) setCurrTicket(ctx context.Context, q pgsql.Querier, userID, ticketID int) error {
	_, err := q.ExecContext(ctx, setCurrTicketSQL, userID, ticketID)
	return err
}

func (p *PgRepo) SetCurrTicket(ctx context.Context, userID, ticketID int) error {
	return p.setCurrTicket(ctx, p.conn, userID, ticketID)
}

func (p *PgRepo) SetCurrTicketTx(ctx context.Context, tx repository.Transaction, userID, ticketID int) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := p.setCurrTicket(ctx, sqlTx, userID, ticketID)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}
