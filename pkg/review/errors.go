package review

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	// Error codes

	InvalidInputError = errcode.InvalidInputError
	UnknownError      = errcode.UnknownError

	ReviewAlreadySubmittedError = errcode.New(
		"review_already_submitted",
		errors.New("ticket review is already submitted"),
	)
)
