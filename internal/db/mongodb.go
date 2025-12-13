package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var DbFileName string = "db.go"

func Connect(ctx context.Context, uri string) (client *mongo.Client, err error) {
	return
}

func Close(ctx context.Context, client *mongo.Client) {
}
