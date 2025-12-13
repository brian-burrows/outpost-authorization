package auth

import "context"

var ServiceFileName string = "service.go"

type Token struct {
}

type TokenService interface {
	Authenticate(ctx context.Context, email, password string) (token string, err error)
	ValidateToken(ctx context.Context, token string) (claims interface{}, err error)
	IssueRefreshToken(ctx context.Context, refreshToken string) (token string, err error)
}

type BcryptTokenService struct {
	repo      UserRepository
	jwtSecret string
}

func NewBcryptTokenService(r UserRepository, secret string) TokenService {
	return &BcryptTokenService{repo: r, jwtSecret: secret}
}

func (t *BcryptTokenService) Authenticate(ctx context.Context, email, password string) (token string, err error) {
	// AWS provides `Amazon Time Sync Service` that should make `time.Now().Unix()` reliable
	token = "hello"
	err = nil
	return
}

type JwtClaims struct {
	iss      string
	sub      string
	aud      string
	exp      float64
	nbf      float64
	iat      float64
	jti      string
	user_id  int
	roles    []string
	username string
}

func (t *BcryptTokenService) ValidateToken(ctx context.Context, token string) (claims interface{}, err error) {
	// Validate the JWT using the payload and the jwtSecret
	// Ensure that the time is between `nbf` and `exp` times, based on `iat`
	// Add leeway to the times to ensure clock drift between servers is acceptable
	// Do something with the `jti` to prevent replay attacks?
	// If the token is valid, then build the JwtClaims data structure
	// Otherwise, return an error
	return
}

func (t *BcryptTokenService) IssueRefreshToken(ctx context.Context, refresh_token string) (token string, err error) {
	// Decode the refresh_token and extract the user information
	// fetch refresh token for the provided username from the database
	// validate that it matches the refresh_token provided
	// validate that the refresh token itself is not expired
	// If everything is okay, create a new token
	// TODO: build a private helper function that doesn't rely on `password` to do so
	return
}
