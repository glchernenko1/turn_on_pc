package db

import (
	"context"
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
		return u, fmt.Errorf("failed to find user by id: %s. err: %v", id, err)
	}
	err = result.Decode(&u)
	if err != nil {
		return u, fmt.Errorf("failed to decode user(id:%s) frone db. err: %v", id, err)
	}
	return u, nil

}

func (d *db) Update(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (d *db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
