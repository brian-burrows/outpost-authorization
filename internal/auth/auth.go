package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
)

type User struct {
	ID    string
	Email string
}

type Role struct{}
type Permission struct{}

var RegisteredEmails = make(map[string]bool)

func CreateUser(email string, providerType string, providerKey string, credential string) (user *User, err error) {
	if RegisteredEmails[email] {
		return &User{}, errors.New("duplicate email address found")
	}
	b := make([]byte, 8)
	rand.Read(b)
	id := fmt.Sprintf("%x", b)
	user = &User{
		ID:    id,
		Email: email,
	}
	RegisteredEmails[email] = true
	return
}
