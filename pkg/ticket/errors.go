package ticket

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	// Error codes

	InvalidInputError = errcode.InvalidInputError
	UnknownError      = errcode.UnknownError

	TicketAlreadyEndedError = errcode.New(
		"ticket_already_ended",
		errors.New("ticket is already ended"),
	)
)
