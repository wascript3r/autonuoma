package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validate struct {
	govalidate *validator.Validate
}

func New() *Validate {
	goV := validator.New()

	r := newRules()
	r.attachTo(goV)

	return &Validate{goV}
}

func (v *Validate) RawRequest(s interface{}) error {
	return v.govalidate.Struct(s)
}
