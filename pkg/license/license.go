package license

import (
	"io"
	"time"

	"github.com/wascript3r/autonuoma/pkg/user"
)

// Confirm, Reject

type ChangeStatusReq struct {
	LicenseID int `json:"licenseID" validate:"required"`
}

// GetAll

type LicenseListInfo struct {
	ID         int                     `json:"id"`
	Number     string                  `json:"number"`
	Client     *user.UserSensitiveInfo `json:"client"`
	Expiration time.Time               `json:"expiration"`
}

type GetAllRes struct {
	Licenses []*LicenseListInfo `json:"licenses"`
}

// GetPhotos

type GetPhotosReq struct {
	LicenseID int `json:"licenseID" validate:"required"`
}

type PhotoListInfo struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type GetPhotosRes struct {
	Photos []*PhotoListInfo `json:"photos"`
}

// UploadLicense

const (
	LicenseKey               = "license"
	LicenseExpirationDateKey = "expirationDate"
	LicenseNumberKey         = "number"
)

type UploadReq struct {
	File                  io.Reader
	LicenseExpirationDate time.Time
	LicenseNumber         string
	Uid                   int
}

type UploadRes struct {
	Filename      string `json:"filename"`
	LicenseStatus string `json:"license"`
}
