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
	STUDIOS  = "studios"
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

func (c *MongoClient) CreateStudio(s models.Studio) error {
	_, err := c.cli.Database(DATABASE).Collection(STUDIOS).InsertOne(context.TODO(), s)
	if err != nil {
		log.Error("insert fail:", err)
		return err
	}
	return nil
}

func (c *MongoClient) ListStudios() ([]models.Studio, error) {
	var studios []models.Studio
	cursor, err := c.cli.Database(DATABASE).Collection(STUDIOS).Find(context.TODO(), bson.D{})
	if err = cursor.All(context.TODO(), &studios); err != nil {
		log.Error("list fail:", err)
		return studios, err
	}

	return studios, nil
}

func (c *MongoClient) DeleteStudio(id string) error {
	_, err := c.cli.Database(DATABASE).Collection(STUDIOS).DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (c *MongoClient) JoinPlayer(studioID string, playerID string) error {
	filter := bson.M{"id": studioID}
	update := bson.M{
		"$addToSet": bson.M{"players": playerID},
	}

	// TODO: 중복 아이디 처리
	result := c.cli.Database(DATABASE).Collection(STUDIOS).FindOneAndUpdate(context.TODO(), filter, update, nil)
	if result.Err() != nil {
		log.Error(result.Err())
		return result.Err()
	}
	return nil
}

func (c *MongoClient) LeavePlayer(studioID string, playerID string) error {
	filter := bson.M{"id": studioID}
	update := bson.M{
		"$pull": bson.M{"players": playerID},
	}

	// TODO: not found id 처리
	result := c.cli.Database(DATABASE).Collection(STUDIOS).FindOneAndUpdate(context.TODO(), filter, update, nil)
	if result.Err() != nil {
		log.Error(result.Err())
		return result.Err()
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
