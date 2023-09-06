package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBConnection() *mongo.Client {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}
	return client
}

func MigrateCollections() {
	dbName := "mongo"
	ctx := context.TODO()

	client := NewDBConnection()
	defer client.Disconnect(ctx)

	err := client.Database(dbName).CreateCollection(ctx, dbName)
	if err != nil {
		return
	}
	err = client.Database(dbName).CreateCollection(ctx, "users")
	if err != nil {
		fmt.Println(err)
		return
	}
}
