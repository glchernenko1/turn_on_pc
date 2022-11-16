package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	var mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if username != "" && password != "" {
		if authDB == "" {
			authDB = database
		}
		clientOptions = clientOptions.SetAuth(options.Credential{
			AuthMechanism: authDB,
			Username:      username,
			Password:      password,
		})
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongoDB: %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error pinging mongoDB: %v", err)
	}
	return client.Database(database), nil
}
