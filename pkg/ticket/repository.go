package ticket

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, ts *domain.Ticket) error
	InsertTx(ctx context.Context, tx repository.Transaction, ts *domain.Ticket) error

	GetLastActiveTicketID(ctx context.Context, clientID int) (int, error)
	GetLastActiveTicketIDTx(ctx context.Context, tx repository.Transaction, clientID int) (int, error)

	IsLastTicketEnded(ctx context.Context, clientID int) (bool, error)
	IsLastTicketEndedTx(ctx context.Context, tx repository.Transaction, clientID int) (bool, error)
}
