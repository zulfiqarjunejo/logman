package clients

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientModel interface {
	FindClientById(id string) (Client, error)
}

type MongoClientModel struct {
	Mongo *mongo.Client
}

func NewMongoClientModel(mongo *mongo.Client) MongoClientModel {
	return MongoClientModel{
		Mongo: mongo,
	}
}

func (clientModel *MongoClientModel) FindClientById(id string) (Client, error) {
	clientsCollection := clientModel.Mongo.Database("logs").Collection("clients")

	var client Client

	// TODO: Learn this syntax in details
	filter := bson.D{{Key: "client_id", Value: id}}
	err := clientsCollection.FindOne(context.Background(), filter).Decode(&client)
	if err != nil {
		return client, err
	}

	return client, nil
}
