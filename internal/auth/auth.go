package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
)

type User struct {
	ID          string
	Email       string
	Credentials map[string]string
}

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidProvider = errors.New("invalid provider key")
	randReader         = rand.Read
)

type AuthorizationService struct {
	registry map[string]*User
}

func NewAuthorizationService() *AuthorizationService {
	return &AuthorizationService{registry: make(map[string]*User)}
}

type ErrDuplicateField struct {
	Field string
	Value string
}

func (e *ErrDuplicateField) Error() string {
	return fmt.Sprintf("duplicate %s found: %s", e.Field, e.Value)
}

func makeRandomIdentifier() (id string, err error) {
	b := make([]byte, 8)
	_, err = randReader(b)
	if err != nil {
		return "", err
	}
	id = fmt.Sprintf("%x", b)
	return
}

func (auth *AuthorizationService) CreateUser(email string, providerType string, providerKey string, credential string) (user *User, err error) {
	validations := []struct {
		isValid bool
		err     error
	}{
		{len(email) >= 6 && strings.Contains(email, "@"), ErrInvalidEmail},
		{len(providerKey) >= 1, ErrInvalidProvider},
	}
	for _, v := range validations {
		if !v.isValid {
			return nil, fmt.Errorf("creating user: %w", v.err)
		}
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
		_, ok := auth.registry[field.Value]
		if ok {
			return nil, &ErrDuplicateField{Field: field.Label, Value: field.Value}
		}
	}
	id, err := makeRandomIdentifier()
	if err != nil {
		return nil, err
	}
	user = &User{
		ID:          id,
		Email:       email,
		Credentials: map[string]string{providerKeyPath: credential},
	}
	// register the user
	for _, field := range checkFields {
		auth.registry[field.Value] = user
	}
	return
}

func (auth *AuthorizationService) GetUser(email string) (*User, error) {
	registrationKey := fmt.Sprintf("reg:%s", email)
	if user, ok := auth.registry[registrationKey]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("unable to fetch target user using email %s", email)
}

func (auth *AuthorizationService) Login(providerType, providerKey, credential string) (*User, error) {
	providerKeyPath := fmt.Sprintf("providerInfo:%s:%s", providerType, providerKey)
	user, ok := auth.registry[providerKeyPath]
	if ok && user.Credentials[providerKeyPath] == credential {
		return user, nil
	}
	return nil, fmt.Errorf("unable to authenticate, invalid credentials")
}
