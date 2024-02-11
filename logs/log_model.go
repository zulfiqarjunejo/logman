package logs

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogModel struct {
	Mongo *mongo.Client
}

func NewLogModel(mongo *mongo.Client) LogModel {
	return LogModel{
		Mongo: mongo,
	}
}

func (logModel *LogModel) Create(log Log) error {
	// TODO: Revise the code after learning MongoDB with Golang.
	logsCollection := logModel.Mongo.Database("logs").Collection("logs")
	_, err := logsCollection.InsertOne(context.Background(), log)
	return err
}

func (logModel *LogModel) GetAll() ([]Log, error) {
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
