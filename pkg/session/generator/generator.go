package generator

import (
	"crypto/rand"

	"github.com/wascript3r/gocipher/encoder"
	"github.com/wascript3r/gocipher/sha256"
)

type Generator struct{}

func New() Generator {
	return Generator{}
}

func (g Generator) GenerateID() (string, error) {
	bs := make([]byte, 64)

	_, err := rand.Read(bs)
	if err != nil {
		return "", err
	}
	hex := encoder.HexEncode(sha256.Compute(bs))

	return string(hex), nil
}
