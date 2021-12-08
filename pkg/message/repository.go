package message

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, ms *domain.Message) (*domain.MessageFull, error)
	InsertTx(ctx context.Context, tx repository.Transaction, ms *domain.Message) (*domain.MessageFull, error)

	GetByTicket(ctx context.Context, ticketID int) ([]*domain.MessageFull, error)
	GetByTicketTx(ctx context.Context, tx repository.Transaction, ticketID int) ([]*domain.MessageFull, error)
}
