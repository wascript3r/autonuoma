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
		"u_email":    "lte=200,email",
		"u_password": "gte=8,lte=100",
	}

	for k, v := range aliases {
		goV.RegisterAlias(k, v)
	}
}
