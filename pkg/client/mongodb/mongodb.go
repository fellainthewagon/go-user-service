package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, mongodbURL, host, port, username, password, database, authDB string) (*mongo.Database, error) {
	var isAuth bool

	if mongodbURL == "" {
		if username == "" && password == "" {
			mongodbURL = fmt.Sprintf("mongodb://%s:%s", host, port)
		} else {
			mongodbURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
			isAuth = true
		}
	}

	clientOptions := options.Client().ApplyURI(mongodbURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to mongodb: %v\n", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Failed to ping to mongodb: %v\n", err)
	}

	return client.Database(database), nil
}
