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
		"m_message": "lte=100",
	}

	for k, v := range aliases {
		goV.RegisterAlias(k, v)
	}
}
