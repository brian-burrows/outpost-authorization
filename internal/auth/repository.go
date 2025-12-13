package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var RepositoryFileName string = "repository.go"

type User struct{}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (user User, err error)
	CreateUser(ctx context.Context, user User) (err error)
	GetUserRefreshToken(ctx context.Context, user User) (token string)
}

type MongoDbUserRepository struct {
	connectionPool *mongo.Client
}

func (u *MongoDbUserRepository) GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	return
}
func (u *MongoDbUserRepository) CreateUser(ctx context.Context, user User) (err error) { return }
func (u *MongoDbUserRepository) GetUserRefreshToken(ctx context.Context, user User) (token string) {
	return
}

func NewMongoDbUserRepository(dbClient *mongo.Client) (client UserRepository) {
	client = &MongoDbUserRepository{connectionPool: dbClient}
	return
}
