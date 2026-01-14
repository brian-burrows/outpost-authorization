package auth

import "time"

func NewCredentials() {}

type Credentials interface {
	IsValid(credential string) bool
}

type passwordCredentials struct {
	HashedPassword string
}

func (pc passwordCredentials) IsValid(credential string) bool {
	// todo: add password hashing
	return pc.HashedPassword == credential
}

type oAuthCredentials struct {
	accessToken  string
	refreshToken string
	expiry       time.Time
	tokenType    string
}

func (pc oAuthCredentials) IsValid(credential string) bool {
	// todo: add complex logic here
	return pc.accessToken == credential
}
