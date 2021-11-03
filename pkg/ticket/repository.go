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

	GetCurrTicketID(ctx context.Context, clientID int) (int, error)
	GetCurrTicketIDTx(ctx context.Context, tx repository.Transaction, clientID int) (int, error)

	IsCurrTicketEnded(ctx context.Context, clientID int) (bool, error)
	IsCurrTicketEndedTx(ctx context.Context, tx repository.Transaction, clientID int) (bool, error)
}
