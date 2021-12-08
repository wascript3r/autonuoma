package license

// Confirm, Reject

type ChangeStatusReq struct {
	LicenseID int `json:"licenseID" validate:"required"`
}
