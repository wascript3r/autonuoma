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

	TicketNotEndedError = errcode.New(
		"ticket_not_ended",
		errors.New("ticket is not ended"),
	)

	TicketAlreadyAcceptedError = errcode.New(
		"ticket_already_accepted",
		errors.New("ticket is already accepted"),
	)

	TicketNotAcceptedError = errcode.New(
		"ticket_not_accepted",
		errors.New("ticket is not accepted"),
	)

	TicketStillActiveError = errcode.New(
		"current_ticket_still_active",
		errors.New("current ticket is still active"),
	)

	TicketNotFoundError = errcode.New(
		"ticket_not_found",
		errors.New("ticket not found"),
	)

	TicketNotOwnedError = errcode.New(
		"ticket_not_owned",
		errors.New("ticket is not owned by you"),
	)
)
