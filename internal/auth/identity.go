package auth

import (
	"fmt"
	"strings"
)

func NewIdentity(providerType, providerKey string) Identity {
	base := BaseIdentity{
		providerType: providerType,
		providerKey:  providerKey,
	}
	var identity Identity
	if providerType == "email" {
		identity = EmailIdentity{BaseIdentity: base}
	} else {
		identity = DefaultIdentity{BaseIdentity: base}
	}
	return identity
}

type BaseIdentity struct {
	providerType string
	providerKey  string
}

func (b BaseIdentity) ProviderType() string { return b.providerType }
func (b BaseIdentity) ProviderKey() string  { return b.providerKey }

type Identity interface {
	IdentityKey() (string, error)
	ProviderType() string
	ProviderKey() string
}

func RegistryKey(id Identity) string {
	return fmt.Sprintf("providerInfo:%s:%s", id.ProviderType(), id.ProviderKey())
}

type DefaultIdentity struct {
	BaseIdentity
}

func (identity DefaultIdentity) IdentityKey() (string, error) {
	if len(identity.providerKey) < 1 {
		return "", ErrInvalidProvider
	}
	return RegistryKey(identity), nil
}

type EmailIdentity struct {
	BaseIdentity
}

func (identity EmailIdentity) IdentityKey() (string, error) {
	length := len(identity.providerKey)
	minLength := 2
	if length < minLength {
		return "", ErrInvalidProvider
	}
	requiredElements := []string{"@"}
	for _, element := range requiredElements {
		if !strings.Contains(identity.providerKey, element) {
			return "", ErrInvalidProvider
		}
	}
	return RegistryKey(identity), nil
}

type PhoneNumberIdentity struct {
	BaseIdentity
}

func (identity PhoneNumberIdentity) IdentityKey() (string, error) {
	return RegistryKey(identity), nil
}
