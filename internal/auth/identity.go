package auth

import (
	"fmt"
	"strings"
)

type Identity struct {
	providerType string
	providerKey  string
}

func (identity *Identity) minKeyLength() int {
	if identity.providerType == "email" {
		return 6
	}
	return 1
}

func (identity *Identity) requiredKeyCharacters() []string {
	if identity.providerType == "email" {
		return []string{"@"}
	}
	return []string{}
}

func (identity *Identity) IdentityKey() (string, error) {
	length := len(identity.providerKey)
	minLength := identity.minKeyLength()
	if length < minLength {
		return "", ErrInvalidProvider
	}
	requiredElements := identity.requiredKeyCharacters()
	for _, element := range requiredElements {
		if !strings.Contains(identity.providerKey, element) {
			return "", ErrInvalidProvider
		}
	}
	return fmt.Sprintf("providerInfo:%s:%s", identity.providerType, identity.providerKey), nil
}

type EmailIdentity struct{}
type PhoneNumberIdentity struct{}
