package mongodb

import (
	"context"
	"fmt"
	"rest-api/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, cfg config.MongoDBConfig) (*mongo.Client, error) {
	var isAuth bool
	mongodbURL := cfg.AtlasURI

	if mongodbURL == "" {
		if cfg.Username == "" && cfg.Password == "" {
			mongodbURL = fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
		} else {
			mongodbURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
			isAuth = true
		}
	}

	clientOptions := options.Client().ApplyURI(mongodbURL)
	if isAuth {
		if cfg.AuthDB == "" {
			cfg.AuthDB = cfg.Database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: cfg.AuthDB,
			Username:   cfg.Username,
			Password:   cfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to mongodb: %v\n", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Failed to ping to mongodb: %v\n", err)
	}

	return client, nil
}
