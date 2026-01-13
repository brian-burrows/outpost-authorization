package auth

import (
	"crypto/rand"
	"fmt"
)

type UuidInterface interface {
	Generate() (string, error)
}

type RandomIdentifierGenerator struct{}

func (g *RandomIdentifierGenerator) Generate() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b) // Uses the real crypto/rand
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
