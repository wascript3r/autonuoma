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
