package message

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, ms *domain.Message) error
	InsertTx(ctx context.Context, tx repository.Transaction, ms *domain.Message) error
}
