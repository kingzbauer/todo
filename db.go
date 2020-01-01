package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultClient *mongo.Client
	db            *mongo.Database
)

func GetDatabase() (*mongo.Database, error) {
	settings := FromEnvironment()
	var err error

	defaultClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(settings.DBURI))
	if err != nil {
		return nil, err
	}

	if err := defaultClient.Ping(context.TODO()); err != nil {
		return nil, err
	}

	return defaultClient(settings.DBName), nil
}

func CloseDefaultClient() error {
	if defaultClient != nil {
		return defaultClient.Disconnect(context.TODO())
	}
	return nil
}
