package auth

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func MakeMongoClient() (*mongo.Client, error) {
	uri := "mongodb://mongodb:27017/?maxPoolSize=50&minPoolSize=10&maxIdleTimeMS=30000"
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	return client, err
}

func TestMongoConnection(t *testing.T) {
	_, err := MakeMongoClient()
	if err != nil {
		t.Fatalf("Failed to connect to the client")
	}
}

func TestRepositoryContract(t *testing.T) {
	client := &mongo.Client{}
	repo := MongoUserRepository{client}
	user := &User{
		ID: "user-123",
		Identities: []Identity{
			NewIdentity("email", "test@example.com"),
		},
	}
	if err := repo.Save(user); err != nil {
		t.Fatalf("expected no error on save, got %v", err)
	}
	found, err := repo.Find("email", "test@example.com")
	if err != nil {
		t.Fatalf("expected to find user, got %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, found.ID)
	}
}
