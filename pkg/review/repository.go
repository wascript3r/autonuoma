package review

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, rs *domain.Review) error
	InsertTx(ctx context.Context, tx repository.Transaction, rs *domain.Review) error

	GetByTicketID(ctx context.Context, ticketID int) (*domain.Review, error)
	GetByTicketIDTx(ctx context.Context, tx repository.Transaction, ticketID int) (*domain.Review, error)
}
