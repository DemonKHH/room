package db

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

var client *mongo.Client

func init() {
	ConnectDB()
}

func ConnectDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	client, err = mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Print("Connection failed to MongoDB")
		log.Fatal(err)
	}

	fmt.Print("Connected to MongoDB")
}

func GetMongoClient() *mongo.Client {
	return client
}

func OpenCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	return client.Database("wmt").Collection(CollectionName)
}
