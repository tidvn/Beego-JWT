package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	Users  *mongo.Collection
}

func (db *DB) Close() {
	err := db.Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func Connect() (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {

		log.Fatal(err)
	}
	database := client.Database("DPay_DB")

	users := database.Collection("Users")
	email_unique_index := mongo.IndexModel{
		Keys:    bson.M{"Email": 1},
		Options: options.Index().SetUnique(true).SetName("email_unique_index"),
	}
	username_unique_index := mongo.IndexModel{
		Keys:    bson.M{"UserName": 1},
		Options: options.Index().SetUnique(true).SetName("username_unique_index"),
	}
	_, err = users.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{email_unique_index, username_unique_index})
	return &DB{
		Client: client,
		Users:  users,
	}, nil
}
