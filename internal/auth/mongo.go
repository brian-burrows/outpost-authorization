package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserRepository struct {
	client *mongo.Client
}

func (repo MongoUserRepository) Find(ctx context.Context, pType, pKey string) (*User, error) {
	coll := repo.client.Database("auth").Collection("users")

	// Search inside the array of identities
	filter := bson.M{
		"identities.provider_type": pType,
		"identities.provider_key":  pKey,
	}
	var user User
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (repo MongoUserRepository) Save(ctx context.Context, user *User) error {
	coll := repo.client.Database("auth").Collection("users")
	opts := options.Replace().SetUpsert(true)
	_, err := coll.ReplaceOne(ctx, bson.M{"_id": user.ID}, user, opts)
	return err
}

func (repo MongoUserRepository) Initialize(ctx context.Context) error {
	coll := repo.client.Database("auth").Collection("users")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "identities.provider_type", Value: 1},
			{Key: "identities.provider_key", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("unique_identity"),
	}
	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	return err
}
