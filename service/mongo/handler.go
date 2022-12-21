package db

import (
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const dataBaseName = "wmt"

func Insert(client *mongo.Client, collectionName string, document interface{}) error {
	_, err := client.Database(dataBaseName).Collection(collectionName).InsertOne(context.TODO(), document)
	return err
}

func Find(client *mongo.Client, collectionName string, filter interface{}) ([]byte, error) {
	var raw []byte
	var err error
	collection := client.Database(dataBaseName).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	raw, err = json.Marshal(results)
	return raw, err
}

func Delete(client *mongo.Client, collectionName string, filter interface{}) error {
	_, err := client.Database(dataBaseName).Collection(collectionName).DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func Update(client *mongo.Client, collectionName string, filter bson.D, document interface{}) error {
	_, err := client.Database(dataBaseName).Collection(collectionName).UpdateOne(context.TODO(), filter, document)
	return err
}
