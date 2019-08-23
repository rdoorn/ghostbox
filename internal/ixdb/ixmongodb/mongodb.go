package ixmongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	//db
	database *mongo.Database
}

func New(url, database string) (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	client.Connect(nil)
	return &MongoDB{
		database: client.Database(database),
	}, nil
}
