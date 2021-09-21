package user

import "context"

type Validate interface {
	RawRequest(s interface{}) error
	EmailUniqueness(ctx context.Context, email string) error
}
