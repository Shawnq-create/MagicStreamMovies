package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set in .env file")
	}

	fmt.Println("Connecting to MongoDB at", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil
	}

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	err := godotenv.Load(".env") // Load environment variables from .env file if no mistakes err returns nil
	if err != nil {
		log.Println("Error loading .env file")
	}
	databaseName := os.Getenv("DATABASE_NAME")
	fmt.Println("Using database:", databaseName)
	collection := Client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		log.Fatalf("Failed to open collection: %s", collectionName)
	}
	return collection
}
