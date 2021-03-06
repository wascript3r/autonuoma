package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/wascript3r/autonuoma/pkg/message"
)

type Validate struct {
	govalidate *validator.Validate
}

func New(mv message.Validate) *Validate {
	goV := validator.New()
	mv.AttachRules(goV)
	return &Validate{goV}
}

func (v *Validate) RawRequest(s interface{}) error {
	return v.govalidate.Struct(s)
}
