package mongo

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	cli *mongo.Client
}

func NewMongoConn() *MongoClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return &MongoClient{client}
}

func (c *MongoClient) Close() {
	if err := c.cli.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
