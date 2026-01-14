package auth

import (
	"fmt"
	"strings"
)

type baseIdentity struct {
	providerType string
	providerKey  string
	credentials  Credentials
}

type IdentityOption func(*baseIdentity)

func WithCredentials(credentials Credentials) IdentityOption {
	return func(b *baseIdentity) {
		b.credentials = credentials
	}
}

func NewIdentity(providerType, providerKey string, identityOpts ...IdentityOption) Identity {
	base := baseIdentity{
		providerType: providerType,
		providerKey:  providerKey,
		credentials:  NoCredentials{},
	}
	for _, opt := range identityOpts {
		opt(&base)
	}
	var identity Identity
	if providerType == "email" {
		identity = emailIdentity{baseIdentity: base}
	} else if providerType == "phone" {
		identity = phoneNumberIdentity{baseIdentity: base}
	} else {
		identity = defaultIdentity{baseIdentity: base}
	}
	return identity
}

func (b baseIdentity) ProviderType() string     { return b.providerType }
func (b baseIdentity) ProviderKey() string      { return b.providerKey }
func (b baseIdentity) Credentials() Credentials { return b.credentials }
func (b baseIdentity) Matches(providerType, providerKey string) bool {
	return providerType == b.ProviderType() && providerKey == b.ProviderKey()
}
func (b baseIdentity) Validate(attempt string) bool { return b.credentials.IsValid(attempt) }

type Identity interface {
	IdentityKey() (string, error)
	ProviderType() string
	ProviderKey() string
	Credentials() Credentials
	Matches(providerType, providerKey string) bool
	Validate(attempt string) bool
}

func RegistryKey(id Identity) string {
	return fmt.Sprintf("providerInfo:%s:%s", id.ProviderType(), id.ProviderKey())
}

type defaultIdentity struct {
	baseIdentity
}

func (identity defaultIdentity) IdentityKey() (string, error) {
	if len(identity.providerKey) < 1 {
		return "", ErrInvalidProvider
	}
	return RegistryKey(identity), nil
}

type emailIdentity struct {
	baseIdentity
}

func (identity emailIdentity) IdentityKey() (string, error) {
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

type phoneNumberIdentity struct {
	baseIdentity
}

func (id phoneNumberIdentity) IdentityKey() (string, error) {
	if !strings.HasPrefix(id.providerKey, "+") {
		return "", fmt.Errorf("invalid phone: must start with +")
	}
	return RegistryKey(id), nil
}
