package usecase

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/license"
	"github.com/wascript3r/autonuoma/pkg/user"
)

type Usecase struct {
	licenseRepo license.Repository
	ctxTimeout  time.Duration

	validate license.Validate

	licensePrefix string
	licenseDir    string
}

func New(lr license.Repository, t time.Duration, v license.Validate, prefix, dir string) *Usecase {
	return &Usecase{
		licenseRepo: lr,
		ctxTimeout:  t,

		validate: v,

		licensePrefix: prefix,
		licenseDir:    dir,
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

func (u *Usecase) GetAllUnconfirmed(ctx context.Context) (*license.GetAllRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ls, err := u.licenseRepo.GetAllUnconfirmed(c)
	if err != nil {
		return nil, err
	}

	licenses := make([]*license.LicenseListInfo, len(ls))
	for i, l := range ls {
		licenses[i] = &license.LicenseListInfo{
			ID:     l.ID,
			Number: l.Number,
			Client: &user.UserSensitiveInfo{
				UserInfo: &user.UserInfo{
					ID:        l.ClientMeta.ID,
					FirstName: l.ClientMeta.FirstName,
					LastName:  l.ClientMeta.LastName,
				},
				PIN: l.ClientMeta.PIN,
			},
			Expiration: l.Expiration,
		}
	}

	return &license.GetAllRes{
		Licenses: licenses,
	}, nil
}

func (u *Usecase) GetPhotos(ctx context.Context, req *license.GetPhotosReq) (*license.GetPhotosRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	status, err := u.licenseRepo.GetStatus(c, req.LicenseID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, license.LicenseNotFoundError
		}
		return nil, err
	}

	if status != domain.SubmittedLicenseStatus {
		return nil, license.LicenseAlreadyProcessedError
	}

	ps, err := u.licenseRepo.GetPhotos(c, req.LicenseID)
	if err != nil {
		return nil, err
	}

	photos := make([]*license.PhotoListInfo, len(ps))
	for i, p := range ps {
		photos[i] = &license.PhotoListInfo{
			ID:  p.ID,
			URL: p.URL,
		}
	}

	return &license.GetPhotosRes{
		Photos: photos,
	}, nil
}

func (u *Usecase) Upload(ctx context.Context, req *license.UploadReq) (*license.UploadRes, error) {
	tempFile, err := ioutil.TempFile(u.licenseDir, fmt.Sprintf("%s-*", u.licensePrefix))
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, req.File)
	if err != nil {
		return nil, err
	}

	status, err := u.licenseRepo.UploadLicense(ctx, req.Uid, req.LicenseExpirationDate, req.LicenseNumber, filepath.Base(tempFile.Name()))
	if err != nil {
		return nil, err
	}

	return &license.UploadRes{
		Filename:      tempFile.Name(),
		LicenseStatus: status,
	}, nil
}
