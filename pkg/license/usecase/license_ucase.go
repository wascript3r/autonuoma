package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/license"
)

type Usecase struct {
	licenseRepo license.Repository
	ctxTimeout  time.Duration

	validate license.Validate
}

func New(lr license.Repository, t time.Duration, v license.Validate) *Usecase {
	return &Usecase{
		licenseRepo: lr,
		ctxTimeout:  t,

		validate: v,
	}
}

func (u *Usecase) changeStatus(ctx context.Context, newStatus domain.LicenseStatus, req *license.ChangeStatusReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return license.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.licenseRepo.NewTx(c)
	if err != nil {
		return err
	}

	status, err := u.licenseRepo.GetStatusTx(c, tx, req.LicenseID)
	if err != nil {
		if err == domain.ErrNotFound {
			return license.LicenseNotFoundError
		}
		return err
	}

	if status != domain.SubmittedLicenseStatus {
		return license.LicenseAlreadyProcessedError
	}

	err = u.licenseRepo.SetStatusTx(c, tx, req.LicenseID, newStatus)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Confirm(ctx context.Context, req *license.ChangeStatusReq) error {
	return u.changeStatus(ctx, domain.ConfirmedLicenseStatus, req)
}

func (u *Usecase) Reject(ctx context.Context, req *license.ChangeStatusReq) error {
	return u.changeStatus(ctx, domain.RejectedLicenseStatus, req)
}
