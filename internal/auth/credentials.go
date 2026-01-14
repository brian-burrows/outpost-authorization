package auth

import "time"

func NewCredentials() {}

type Credentials interface {
	IsValid(attempt string) bool
}

type NoCredentials struct{}

func (pc NoCredentials) IsValid(attempt string) bool {
	return false
}

type PasswordCredentials struct {
	hashedPassword string
}

func (pc PasswordCredentials) IsValid(attempt string) bool {
	// todo: add password hashing
	return pc.hashedPassword == attempt
}

type OauthCredentials struct {
	accessToken  string
	refreshToken string
	expiry       time.Time
	tokenType    string
}

func (pc OauthCredentials) IsValid(attempt string) bool {
	// todo, ensure that the accessToken hasn't expired. Try to refresh?
	return pc.accessToken == attempt
}
