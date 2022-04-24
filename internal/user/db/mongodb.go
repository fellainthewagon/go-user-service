package db

import (
	"context"
	"errors"
	"fmt"
	"rest-api/internal/apperror"
	"rest-api/internal/user"
	"rest-api/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("Create user...")

	d.logger.Infoln(user)
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("MongoDB: Failed to insert: %v", err)
	}

	oID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("Failed ObjectID to hex conversion: %v", err)
	}

	return oID.Hex(), nil
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	d.logger.Debug("FindOne user...")

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("Failed hex (id: %s) to ObjectID conversion: %v", id, err)
	}

	result := d.collection.FindOne(ctx, bson.M{"_id": oID})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.NotFoundError
		}
		return u, fmt.Errorf("MongoDB: Failed find user (id: %s).\n Error: %v", id, err)
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("Failed decoding user (id: %s).\n Error: %v", id, err)
	}

	return u, err
}

func (d *db) FindAll(ctx context.Context) (users []user.User, err error) {
	d.logger.Debug("FindAll users...")

	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return users, fmt.Errorf("MongoDB: Failed get cursor: %v", err)
	}

	if err = cursor.All(ctx, &users); err != nil {
		return users, fmt.Errorf("MongoDB: Failed iteration of cursor: %v", err)
	}

	return users, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	d.logger.Debug("Update user...")

	oID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("Failed hex (id: %s) to ObjectID conversion: %v", user.ID, err)
	}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to marshal user (id: %s). Error: %v", user.ID, err)
	}

	var updateUserData bson.M
	err = bson.Unmarshal(userBytes, &updateUserData)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal userBytes: %v", err)
	}

	delete(updateUserData, "_id")
	filter := bson.M{"_id": oID}
	update := bson.M{"$set": updateUserData}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("MongoDB: Failed update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return apperror.NotFoundError
	}

	d.logger.Tracef("MatchedCount: %d, modified: %d", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	d.logger.Debug("Delete user...")

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Failed hex (id: %s) to ObjectID conversion: %v", id, err)
	}

	result, err := d.collection.DeleteOne(ctx, bson.M{"_id": oID})
	if err != nil {
		return fmt.Errorf("MongoDB: Failed delete user (id: %s). Error: %v", id, err)
	}

	if result.DeletedCount == 0 {
		return apperror.NotFoundError
	}

	d.logger.Tracef("DeletedCount: %d.", result.DeletedCount)
	return nil
}
