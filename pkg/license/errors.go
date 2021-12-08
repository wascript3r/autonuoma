package license

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	// Error codes

	InvalidInputError = errcode.InvalidInputError
	UnknownError      = errcode.UnknownError

	LicenseNotFoundError = errcode.New(
		"license_not_found",
		errors.New("license not found"),
	)

	LicenseAlreadyProcessedError = errcode.New(
		"license_already_processed",
		errors.New("license is already processed"),
	)
)
