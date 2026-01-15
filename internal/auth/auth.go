package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
)

type ErrDuplicateField struct {
	Field string
	Value string
}

func (e *ErrDuplicateField) Error() string {
	return fmt.Sprintf("duplicate %s found: %s", e.Field, e.Value)
}

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidProvider = errors.New("invalid provider key")
	randReader         = rand.Read
)

type AuthorizationService struct {
	UuidService UuidInterface
	repo        UserRepository
}

func NewAuthorizationService(opts ...func(*AuthorizationService)) *AuthorizationService {
	service := &AuthorizationService{
		UuidService: &RandomIdentifierGenerator{},
		repo:        &InMemoryUserRepository{registry: map[string]*User{}},
	}
	for _, o := range opts {
		o(service)
	}
	return service
}

func WithRepository(repository UserRepository) func(*AuthorizationService) {
	return func(s *AuthorizationService) {
		s.repo = repository
	}
}

func WithUuidService(generator UuidInterface) func(*AuthorizationService) {
	return func(s *AuthorizationService) {
		s.UuidService = generator
	}
}

func (auth *AuthorizationService) CreateUser(
	email,
	providerType,
	providerKey string,
	credential Credentials,

) (user *User, err error) {
	randomId, err := auth.UuidService.Generate()
	if err != nil {
		return nil, err
	}
	identities := []Identity{
		NewIdentity(providerType, providerKey, WithCredentials(credential)),
		NewIdentity("email", email, WithCredentials(credential)),
		NewIdentity("UserId", randomId, WithCredentials(credential)),
	}
	user = &User{
		ID:         randomId,
		Email:      email,
		Identities: identities,
	}
	err = auth.repo.Save(user)
	if err != nil {
		return nil, err
	}
	return
}

func (auth *AuthorizationService) GetUser(email string) (*User, error) {
	user, err := auth.GetUserByIdentity("email", email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (auth *AuthorizationService) GetUserByIdentity(pType, pKey string) (*User, error) {
	return auth.repo.Find(pType, pKey)
}

func (auth *AuthorizationService) findUserById(userId string) (*User, error) {
	return auth.repo.Find("UserId", userId)
}

func (auth *AuthorizationService) Login(pType, pKey string, credential string) (*User, error) {
	if user, err := auth.repo.Find(pType, pKey); err == nil && user.Authenticate(pType, pKey, credential) {
		return user, nil
	}
	return nil, fmt.Errorf("unable to authenticate, invalid credentials")
}

func (auth *AuthorizationService) AddIdentity(userId, pType, pKey string, credential Credentials) error {
	user, err := auth.findUserById(userId)
	if err != nil { // user not found
		return ErrInvalidProvider
	}
	newIdentity := NewIdentity(pType, pKey, WithCredentials(credential))
	if err != nil {
		return err
	}
	updatedUser := user.AddIdentity(newIdentity)
	return auth.repo.Save(updatedUser)
}
