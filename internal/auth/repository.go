package auth

import (
	"context"
	"fmt"
)

func NewDatabase() {}

type UserRepository interface {
	Find(ctx context.Context, pType, pKey string) (*User, error)
	Save(ctx context.Context, user *User) error
}

type InMemoryUserRepository struct {
	registry map[string]*User
}

func (repo InMemoryUserRepository) Find(ctx context.Context, pType, pKey string) (*User, error) {
	key, err := NewIdentity(pType, pKey).IdentityKey()
	if err != nil {
		return nil, err
	}
	if user, ok := repo.registry[key]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found for %s:%s", pType, pKey)
}
func (repo InMemoryUserRepository) Save(ctx context.Context, user *User) error {
	for _, identity := range user.Identities {
		_, err := identity.IdentityKey()
		if err != nil {
			return err
		}
		existing, err := repo.Find(ctx, identity.ProviderType(), identity.ProviderKey())
		if err == nil && existing.ID != user.ID {
			return &ErrDuplicateField{Field: identity.ProviderType(), Value: identity.ProviderKey()}
		}
	}
	for _, identity := range user.Identities {
		key, _ := identity.IdentityKey()
		repo.registry[key] = user
	}
	return nil
}
