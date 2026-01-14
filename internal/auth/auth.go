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

type User struct {
	ID         string
	Email      string
	Identities []Identity
}

func (u *User) Authenticate(providerType, providerKey, attempt string) bool {
	for _, userAlias := range u.Identities {
		if userAlias.Matches(providerType, providerKey) {
			return userAlias.Validate(attempt)
		}
	}
	return false
}

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidProvider = errors.New("invalid provider key")
	randReader         = rand.Read
)

type AuthorizationService struct {
	registry    map[string]*User
	UuidService UuidInterface
}

func NewAuthorizationService(opts ...func(*AuthorizationService)) *AuthorizationService {
	service := &AuthorizationService{
		registry:    make(map[string]*User),
		UuidService: &RandomIdentifierGenerator{},
	}
	for _, o := range opts {
		o(service)
	}
	return service
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
	var keys []string
	for _, identity := range identities {
		key, err := identity.IdentityKey()
		if err != nil {
			return nil, err
		}
		_, ok := auth.registry[key]
		if ok {
			return nil, &ErrDuplicateField{Field: identity.ProviderType(), Value: identity.ProviderKey()}
		}
		keys = append(keys, key)
	}
	user = &User{
		ID:         randomId,
		Email:      email,
		Identities: identities,
	}
	// register the user
	for _, key := range keys {
		auth.registry[key] = user
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
	key, err := NewIdentity(pType, pKey).IdentityKey()
	if err != nil {
		return nil, err
	}
	if user, ok := auth.registry[key]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("Failed to find user by default")
}

func (auth *AuthorizationService) Login(pType, pKey string, credential string) (*User, error) {
	providerKeyPath, err := NewIdentity(pType, pKey).IdentityKey()
	if err != nil {
		return nil, err
	}
	user, ok := auth.registry[providerKeyPath]
	if ok && user.Authenticate(pType, pKey, credential) {
		return user, nil
	}
	return nil, fmt.Errorf("unable to authenticate, invalid credentials")
}

func (auth *AuthorizationService) findUserById(userId string) (*User, error) {
	providerKeyPath, err := NewIdentity("UserId", userId).IdentityKey()
	if err != nil {
		return nil, err
	}
	user, ok := auth.registry[providerKeyPath]
	if ok {
		return user, nil
	}
	return nil, fmt.Errorf("User ID does not exist")
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
	providerKeyPath, err := newIdentity.IdentityKey()
	if err != nil {
		return err
	}
	if otherUser, ok := auth.registry[providerKeyPath]; ok && user.ID != otherUser.ID {
		return fmt.Errorf("Unable to add an additional ")
	}
	user.Identities = append(user.Identities, newIdentity)
	auth.registry[providerKeyPath] = user
	return nil
}
