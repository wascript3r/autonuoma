package ticket

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, ts *domain.Ticket) error
	InsertTx(ctx context.Context, tx repository.Transaction, ts *domain.Ticket) error

	SetAgent(ctx context.Context, id int, agentID int) error
	SetAgentTx(ctx context.Context, tx repository.Transaction, id int, agentID int) error

	SetEnded(ctx context.Context, id int, ended time.Time) error
	SetEndedTx(ctx context.Context, tx repository.Transaction, id int, ended time.Time) error

	SetAgentEnded(ctx context.Context, id int, agentID int, ended time.Time) error
	SetAgentEndedTx(ctx context.Context, tx repository.Transaction, id int, agentID int, ended time.Time) error

	GetLastActiveID(ctx context.Context, clientID int) (int, error)
	GetLastActiveIDTx(ctx context.Context, tx repository.Transaction, clientID int) (int, error)

	GetMeta(ctx context.Context, id int) (*domain.TicketMeta, error)
	GetMetaTx(ctx context.Context, tx repository.Transaction, id int) (*domain.TicketMeta, error)

	GetAll(ctx context.Context) ([]*domain.TicketFull, error)
	GetAllTx(ctx context.Context, tx repository.Transaction) ([]*domain.TicketFull, error)

	GetByUser(ctx context.Context, userID int) ([]*domain.TicketFull, error)
	GetByUserTx(ctx context.Context, tx repository.Transaction, userID int) ([]*domain.TicketFull, error)
}
