package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
)

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

func (auth *AuthorizationService) CreateUser(email string, providerType string, providerKey string, credential string) (user *User, err error) {
	randomId, err := makeRandomIdentifier()
	if err != nil {
		return nil, err
	}
	providers := []struct {
		pType string
		pKey  string
	}{
		{providerType, providerKey},
		{"email", email},
		{"UserId", randomId},
	}
	var keys []string
	for _, p := range providers {
		key, err := auth.identityKey(p.pType, p.pKey)
		if err != nil {
			return nil, err
		}
		_, ok := auth.registry[key]
		if ok {
			return nil, &ErrDuplicateField{Field: p.pType, Value: p.pKey}
		}
		keys = append(keys, key)
	}
	user = &User{
		ID:          randomId,
		Email:       email,
		Credentials: map[string]string{},
	}
	// register the user
	for _, key := range keys {
		user.Credentials[key] = credential
		auth.registry[key] = user
	}
	return
}

func (auth *AuthorizationService) identityKey(pType, pKey string) (string, error) {
	if pType == "email" && len(pKey) >= 6 && strings.Contains(pKey, "@") {
		return fmt.Sprintf("reg:%s", pKey), nil
	} else if pType == "email" {
		// fmt.Printf("Failing on email %s %s", pType, pKey)
		return "", ErrInvalidEmail
	} else if len(pKey) < 1 {
		return "", ErrInvalidProvider
	}
	return fmt.Sprintf("providerInfo:%s:%s", pType, pKey), nil
}

func (auth *AuthorizationService) GetUserByIdentity(pType, pKey string) (*User, error) {
	key, err := auth.identityKey(pType, pKey)
	if err != nil {
		return nil, err
	}
	if user, ok := auth.registry[key]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("Failed to find user by default")
}

func (auth *AuthorizationService) GetUser(email string) (*User, error) {
	user, err := auth.GetUserByIdentity("email", email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (auth *AuthorizationService) Login(providerType, providerKey, credential string) (*User, error) {
	providerKeyPath, err := auth.identityKey(providerType, providerKey)
	if err != nil {
		return nil, err
	}
	user, ok := auth.registry[providerKeyPath]
	if ok && user.Credentials[providerKeyPath] == credential {
		return user, nil
	}
	return nil, fmt.Errorf("unable to authenticate, invalid credentials")
}

func (auth *AuthorizationService) findUserById(userId string) (*User, error) {
	providerKeyPath, err := auth.identityKey("UserId", userId)
	if err != nil {
		return nil, err
	}
	user, ok := auth.registry[providerKeyPath]
	if ok {
		return user, nil
	}
	return nil, fmt.Errorf("User ID does not exist")
}

func (auth *AuthorizationService) AddIdentity(userId, providerType, providerKey, credential string) error {
	user, err := auth.findUserById(userId)
	if err != nil { // user not found
		return ErrInvalidProvider
	}
	providerKeyPath, err := auth.identityKey(providerType, providerKey)
	if err != nil {
		return err
	}
	existingUser, ok := auth.registry[providerKeyPath]
	if ok && existingUser != nil && existingUser.ID != user.ID { // this identity already exists, and the current user doesn't own it
		return nil
	}
	user.Credentials[providerKeyPath] = credential
	auth.registry[providerKeyPath] = user
	return nil
}
