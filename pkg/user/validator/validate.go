package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/wascript3r/autonuoma/pkg/user"
)

type Validate struct {
	govalidate *validator.Validate
	userRepo   user.Repository
}

func New(ur user.Repository) *Validate {
	goV := validator.New()

	r := newRules()
	r.attachTo(goV)

	return &Validate{goV, ur}
}

func (v *Validate) RawRequest(s interface{}) error {
	return v.govalidate.Struct(s)
}

func (v *Validate) EmailUniqueness(ctx context.Context, email string) error {
	exists, err := v.userRepo.EmailExists(ctx, email)
	if err != nil {
		return err
	}

	if exists {
		return user.ErrEmailExists
	}

	return nil
}
