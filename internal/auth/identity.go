package auth

import (
	"fmt"
	"net/mail"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func NewIdentity(
	providerType,
	providerKey string,
	identityOpts ...IdentityOption,
) Identity {
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
func (b baseIdentity) MarshalBSON() ([]byte, error) {
	return bson.Marshal(struct {
		Type string `bson:"provider_type"`
		Key  string `bson:"provider_key"`
		// We can add credentials here later if needed
	}{
		Type: b.providerType,
		Key:  b.providerKey,
	})
}

type Identity interface {
	IdentityKey() (string, error)
	ProviderType() string
	ProviderKey() string
	Credentials() Credentials
	Matches(providerType, providerKey string) bool
	Validate(attempt string) bool
	MarshalBSON() ([]byte, error)
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
	addr, err := mail.ParseAddress(identity.providerKey)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidProvider, err)
	}
	if addr.Address != identity.providerKey {
		return "", fmt.Errorf("%w: email must be in raw format (user@domain.com)", ErrInvalidProvider)
	}
	return RegistryKey(identity), nil
}

type phoneNumberIdentity struct {
	baseIdentity
}

func (identity phoneNumberIdentity) IdentityKey() (string, error) {
	key := identity.providerKey
	if !strings.HasPrefix(key, "+") || len(key) < 8 {
		return "", fmt.Errorf("invalid phone: must be in E.164 format (e.g., +1234567890)")
	}
	digitsOnly := strings.TrimPrefix(key, "+")
	for _, r := range digitsOnly {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("invalid phone: contains non-digit characters")
		}
	}
	return RegistryKey(identity), nil
}
