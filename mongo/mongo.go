package mongo

import (
	"api-server/models"
	"context"
	"errors"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE = "band4all"
	USERS    = "users"
	SESSIONS = "sessions"
	MQTT     = "mqtt"
	TURN     = "turn"
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

func (c *MongoClient) CheckAuth(username string, password string) error {
	var user models.User
	err := c.cli.Database(DATABASE).Collection(USERS).FindOne(context.TODO(), bson.M{"name": username}).Decode(&user)
	if err != nil {
		log.Error(err)
		return err
	}

	if user.Password != password {
		return errors.New("Unauthorized")
	}

	return nil
}

func (c *MongoClient) CreateSession(s models.Session) error {
	_, err := c.cli.Database(DATABASE).Collection(SESSIONS).InsertOne(context.TODO(), s)
	if err != nil {
		log.Error("insert fail:", err)
		return err
	}
	return nil
}

func (c *MongoClient) ListSessions() ([]models.Session, error) {
	var sessions []models.Session
	cursor, err := c.cli.Database(DATABASE).Collection(SESSIONS).Find(context.TODO(), bson.D{})
	if err = cursor.All(context.TODO(), &sessions); err != nil {
		log.Error("list fail:", err)
		return sessions, err
	}

	return sessions, nil
}

func (c *MongoClient) DeleteSession(id string) error {
	_, err := c.cli.Database(DATABASE).Collection(SESSIONS).DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (c *MongoClient) GetMQTT() (models.MQTTServer, error) {
	var mqtt models.MQTTServer
	err := c.cli.Database(DATABASE).Collection(MQTT).FindOne(context.TODO(), bson.M{}).Decode(&mqtt)
	if err != nil {
		log.Error(err)
		return mqtt, err
	}

	return mqtt, nil
}

func (c *MongoClient) GetTURNs() ([]models.TURNServer, error) {
	var turns []models.TURNServer
	cursor, err := c.cli.Database(DATABASE).Collection(TURN).Find(context.TODO(), bson.M{})
	if err != nil {
		log.Error(err)
		return turns, err
	}

	if err = cursor.All(context.TODO(), &turns); err != nil {
		log.Error(err)
		return turns, err
	}

	return turns, nil
}
