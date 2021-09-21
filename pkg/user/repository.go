package user

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	InsertIfNotExists(ctx context.Context, us *domain.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
	GetIDAndPassword(ctx context.Context, email string) (id int, password string, err error)

	DeductBalance(ctx context.Context, id int, value int64) error
	DeductBalanceTx(ctx context.Context, tx repository.Transaction, id int, value int64) error

	AddBalance(ctx context.Context, id int, value int64) error
	AddBalanceTx(ctx context.Context, tx repository.Transaction, id int, value int64) error
}
