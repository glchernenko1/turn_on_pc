package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"turn_on_pc/internal/user"
	"turn_on_pc/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) FiendAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return u, fmt.Errorf("failed to find all users. err: %v", err)
	}
	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed decode. err: %v", err)
	}
	return u, nil
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("Creating user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("faield to create User due to error %v", err)
	}
	d.logger.Debugf("Converting InsertedID to ObjectID")
	oID, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oID.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("faield to convert objectid to hex")
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}
	filter := bson.M{"_id": oID}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find user by id: %s. err: %v", id, err)
	}
	err = result.Decode(&u)
	if err != nil {
		return u, fmt.Errorf("failed to decode user(id:%s) frone db. err: %v", id, err)
	}
	return u, nil

}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectid. hex: %s", user.ID)
	}
	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. err: %v ", err)
	}

	var updateUserObject bson.M
	err = bson.Unmarshal(userBytes, updateUserObject)
	if err != nil {
		return fmt.Errorf("failed to unmarshal userBytes. err: %v", err)
	}
	delete(updateUserObject, "_id")

	update := bson.M{
		"$set": updateUserObject,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update) //todo просмотреть как работает с UpdateByID
	if err != nil {
		return fmt.Errorf("failed to update user query: %v", err)
	}
	if result.MatchedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Matched count: %d and Modified count: %d", result.MatchedCount, result.ModifiedCount)
	return nil

}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}
	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Deleted count: %d", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
