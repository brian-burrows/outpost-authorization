package auth

import (
	"crypto/rand"
	"fmt"
	"strings"
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

func makeRandomIdentifier() (id string, err error) {
	b := make([]byte, 8)
	_, err = rand.Read(b)
	id = fmt.Sprintf("%x", b)
	return
}

func validateProviderKey(providerKey string) (err error) {
	if len(providerKey) < 1 {
		err = fmt.Errorf("duplicate provider key: %s", providerKey)
	}
	return
}
func validateEmail(email string) (err error) {
	minChars := 6
	if len(email) < minChars {
		err = fmt.Errorf("email must have more than %d characters", minChars)
	} else if !strings.Contains(email, "@") {
		err = fmt.Errorf("email must have '@' symbol, received %s", email)
	}
	return
}

func CreateUser(email string, providerType string, providerKey string, credential string) (user *User, err error) {
	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	if err := validateProviderKey(providerKey); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	registrationKey := fmt.Sprintf("reg:%s", email)
	providerKeyPath := fmt.Sprintf("providerInfo:%s:%s", providerType, providerKey)
	// Check if a registered user exists
	checkFields := []struct {
		Label string
		Value string
	}{
		{"email", registrationKey},
		{"providerInfo", providerKeyPath},
	}
	for _, field := range checkFields {
		if RegisteredProviders[field.Value] {
			return nil, &ErrDuplicateField{Field: field.Label, Value: field.Value}
		}
	}
	id, err := makeRandomIdentifier()
	if err != nil {
		return nil, err
	}
	user = &User{
		ID:    id,
		Email: email,
	}
	// register the user
	for _, field := range checkFields {
		RegisteredProviders[field.Value] = true
	}
	return
}
