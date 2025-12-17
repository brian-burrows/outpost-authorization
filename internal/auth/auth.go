package auth

import (
	"crypto/rand"
	"fmt"
)

type User struct {
	ID    string
	Email string
}

var RegisteredProviders = make(map[string]bool)

type ErrDuplicateField struct {
	Field string
	Value string
}

func (e *ErrDuplicateField) Error() string {
	return fmt.Sprintf("duplicate %s found: %s", e.Field, e.Value)
}

func makeRandomIdentifier() (id string) {
	b := make([]byte, 8)
	rand.Read(b)
	id = fmt.Sprintf("%x", b)
	return
}

func CreateUser(email string, providerType string, providerKey string, credential string) (user *User, err error) {
	registrationKey := fmt.Sprintf("reg:%s:%s", email, providerType)
	providerKeyPath := fmt.Sprintf("pkey:%s", providerKey)
	checkFields := []struct {
		Label string
		Value string
	}{
		{"registrationKey", registrationKey},
		{"providerKeyPath", providerKeyPath},
	}
	for _, field := range checkFields {
		if RegisteredProviders[field.Value] {
			return nil, &ErrDuplicateField{Field: field.Label, Value: field.Value}
		}
	}
	id := makeRandomIdentifier()
	user = &User{
		ID:    id,
		Email: email,
	}
	for _, field := range checkFields {
		RegisteredProviders[field.Value] = true
	}
	return
}
