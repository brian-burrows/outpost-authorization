package auth

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongo(ctx context.Context, t *testing.T) (*mongo.Client, func()) {
	t.Helper() // Marks this function as a test helper

	// 1. Start Container
	mongodbContainer, err := mongodb.Run(ctx, "mongo:latest")
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	// 2. Get Connection String
	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// 3. Connect Client
	client, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatalf("failed to connect to mongo: %s", err)
	}

	// 4. Return the client and a cleanup function
	teardown := func() {
		if err := client.Disconnect(ctx); err != nil {
			t.Errorf("failed to disconnect client: %s", err)
		}
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: %s", err)
		}
	}

	return client, teardown
}

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
	ctx := context.Background()
	client, teardown := setupMongo(ctx, t)
	defer teardown()
	repo := MongoUserRepository{client: client}
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

func TestMongoRepository(t *testing.T) {
	ctx := context.Background()
	client, teardown := setupMongo(ctx, t)
	defer teardown()
	repo := MongoUserRepository{client: client}
	user := &User{ID: "user-123"}
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
