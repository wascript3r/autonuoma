package license

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	GetStatus(ctx context.Context, id int) (domain.LicenseStatus, error)
	GetStatusTx(ctx context.Context, tx repository.Transaction, id int) (domain.LicenseStatus, error)

	SetStatus(ctx context.Context, id int, status domain.LicenseStatus) error
	SetStatusTx(ctx context.Context, tx repository.Transaction, id int, status domain.LicenseStatus) error
}
