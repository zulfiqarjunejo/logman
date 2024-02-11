package logs

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogModel interface {
	Create(Log) error
	GetAll() ([]Log, error)
}

type MongoLogModel struct {
	Mongo *mongo.Client
}

func NewMongoLogModel(mongo *mongo.Client) MongoLogModel {
	return MongoLogModel{
		Mongo: mongo,
	}
}

func (logModel *MongoLogModel) Create(log Log) error {
	// TODO: Revise the code after learning MongoDB with Golang.
	logsCollection := logModel.Mongo.Database("logs").Collection("logs")
	_, err := logsCollection.InsertOne(context.Background(), log)
	return err
}

func (logModel *MongoLogModel) GetAll() ([]Log, error) {
	// TODO: Revise the code after learning MongoDB with Golang.
	logsCollection := logModel.Mongo.Database("logs").Collection("logs")

	filter := bson.D{}
	cursor, err := logsCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var logs []Log
	err = cursor.All(context.TODO(), &logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
