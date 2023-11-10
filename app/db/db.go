package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBClient() (*mongo.Client, error) {
	mongoUser     := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost     := os.Getenv("MONGO_HOST")
	
	uri := fmt.Sprintf(`mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority`, mongoUser, mongoPassword, mongoHost)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}
