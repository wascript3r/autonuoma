package validator

import (
	"github.com/go-playground/validator/v10"
)

type rules struct{}

func newRules() rules {
	return rules{}
}

func (r rules) attachTo(goV *validator.Validate) {
	aliases := map[string]string{
		"u_message": "lt=100",
	}

	for k, v := range aliases {
		goV.RegisterAlias(k, v)
	}
}
