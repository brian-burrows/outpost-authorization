package auth

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewDatabase() {}

type UserRepository interface {
	Find(pType, pKey string) (*User, error)
	Save(user *User) error
}

type InMemoryUserRepository struct {
	registry map[string]*User
}

func (repo InMemoryUserRepository) Find(pType, pKey string) (*User, error) {
	key, err := NewIdentity(pType, pKey).IdentityKey()
	if err != nil {
		return nil, err
	}
	if user, ok := repo.registry[key]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found for %s:%s", pType, pKey)
}

func (repo InMemoryUserRepository) Save(user *User) error {
	for _, identity := range user.Identities {
		_, err := identity.IdentityKey()
		if err != nil {
			return err
		}
		existing, err := repo.Find(identity.ProviderType(), identity.ProviderKey())
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

type MongoUserRepository struct {
	client *mongo.Client
}

func (repo MongoUserRepository) Find(pType, pKey string) (*User, error) {
	return &User{ID: "user-123"}, nil
}

func (repo MongoUserRepository) Save(user *User) error {
	return nil
}
