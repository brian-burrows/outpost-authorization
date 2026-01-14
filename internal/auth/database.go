package auth

func NewDatabase() {}

type UserRepository interface {
	CreateUser(identity Identity) error
	GetUser(identity Identity) error
	addIdentity([]Identity) error
}
