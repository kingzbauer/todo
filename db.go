package main

import (
	"context"
	"time"

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

	ctx, _ := context.WithDeadline(context.TODO(), time.Now().Add(time.Second*5))
	defaultClient, err = mongo.Connect(ctx, options.Client().ApplyURI(settings.DBURI))
	if err != nil {
		return nil, err
	}

	if err := defaultClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return defaultClient.Database(settings.DBName), nil
}

func CloseDefaultClient() error {
	if defaultClient != nil {
		return defaultClient.Disconnect(context.TODO())
	}
	return nil
}
