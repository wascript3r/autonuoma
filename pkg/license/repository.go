package license

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	GetStatus(ctx context.Context, id int) (domain.LicenseStatus, error)
	GetStatusTx(ctx context.Context, tx repository.Transaction, id int) (domain.LicenseStatus, error)

	SetStatus(ctx context.Context, id int, status domain.LicenseStatus) error
	SetStatusTx(ctx context.Context, tx repository.Transaction, id int, status domain.LicenseStatus) error

	GetAllUnconfirmed(ctx context.Context) ([]*domain.LicenseFull, error)
	GetAllUnconfirmedTx(ctx context.Context, tx repository.Transaction) ([]*domain.LicenseFull, error)

	GetPhotos(ctx context.Context, licenseID int) ([]*domain.LicensePhoto, error)
	GetPhotosTx(ctx context.Context, tx repository.Transaction, licenseID int) ([]*domain.LicensePhoto, error)

	UploadLicense(ctx context.Context, uid int, expirationDate time.Time, number string, filename string) (string, error)
}
