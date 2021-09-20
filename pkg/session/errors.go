package session

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	ErrCannotLoadSession = errors.New("cannot load session from context")

	// Error codes

	UnknownError = errcode.UnknownError

	NotAuthenticatedError = errcode.New(
		"not_authenticated",
		errors.New("not authenticated"),
	)

	InsufficientPermissionsError = errcode.New(
		"insufficient_permissions",
		errors.New("insufficient permissions"),
	)

	AlreadyAuthenticatedError = errcode.New(
		"already_authenticated",
		errors.New("already authenticated"),
	)

	SessionExpiredError = errcode.New(
		"session_expired",
		errors.New("session is expired"),
	)

	TokenExpiredError = errcode.New(
		"token_expired",
		errors.New("token is expired"),
	)
)
