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
		"r_stars":   "gte=1,lte=5",
		"r_comment": "lt=100",
	}

	for k, v := range aliases {
		goV.RegisterAlias(k, v)
	}
}
