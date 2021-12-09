package domain

import "time"

type LicenseStatus int8

const (
	SubmittedLicenseStatus LicenseStatus = iota + 1
	ConfirmedLicenseStatus
	RejectedLicenseStatus
)

type License struct {
	ID         int
	Number     string
	ClientID   int
	Expiration time.Time
	StatusID   LicenseStatus
}

type LicenseFull struct {
	ID         int
	Number     string
	ClientMeta *UserSensitiveMeta
	Expiration time.Time
	StatusID   LicenseStatus
}

type LicensePhoto struct {
	ID        int
	LicenseID int
	URL       string
}
