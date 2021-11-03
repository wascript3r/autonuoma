package message

import "github.com/go-playground/validator/v10"

type Validate interface {
	RawRequest(s interface{}) error
	AttachRules(goV *validator.Validate)
}
