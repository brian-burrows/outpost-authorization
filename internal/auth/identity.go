package auth

import (
	"fmt"
	"strings"
)

type Identity struct {
	providerType string
	providerKey  string
}

func (*Identity) minKeyLength(pType string) int {
	if pType == "email" {
		return 6
	}
	return 1
}

func (*Identity) requiredKeyCharacters(pType string) []string {
	if pType == "email" {
		return []string{"@"}
	}
	return []string{}
}

func (identity *Identity) IdentityKey(pType, pKey string) (string, error) {
	length := len(pKey)
	minLength := identity.minKeyLength(pType)
	if length < minLength {
		return "", ErrInvalidProvider
	}
	requiredElements := identity.requiredKeyCharacters(pType)
	for _, element := range requiredElements {
		if !strings.Contains(pKey, element) {
			return "", ErrInvalidProvider
		}
	}
	return fmt.Sprintf("providerInfo:%s:%s", pType, pKey), nil
}
