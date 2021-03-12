package mongo

import (
	"api-server/models"
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *MongoClient) databaseStudio() *mongo.Database {
	return c.cli.Database("band4all")
}

func (c *MongoClient) collectionStudio() *mongo.Collection {
	return c.databaseStudio().Collection("studios")
}

func (c *MongoClient) collectionUsers() *mongo.Collection {
	return c.databaseStudio().Collection("users")
}

func (c *MongoClient) collectionMqtt() *mongo.Collection {
	return c.databaseStudio().Collection("mqtt")
}

func (c *MongoClient) collectionTurn() *mongo.Collection {
	return c.databaseStudio().Collection("turn")
}

func (c *MongoClient) CheckAuth(username string, password string) error {
	var user models.User
	err := c.collectionUsers().FindOne(context.TODO(), bson.M{"name": username}).Decode(&user)
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
	_, err := c.collectionStudio().InsertOne(context.TODO(), s)
	if err != nil {
		log.Error("insert fail:", err)
		return err
	}
	return nil
}

// TODO: pagination
func (c *MongoClient) ListStudios(page int, limit int) ([]models.Studio, error) {
	var studios []models.Studio

	opts := options.Find()
	if page != -1 && limit != -1 {
		skip := page * limit
		opts.SetSkip(int64(skip))
	}
	if limit != -1 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := c.collectionStudio().Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Error("find:", err)
		return studios, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &studios); err != nil {
		log.Error("list fail:", err)
		return studios, err
	}

	return studios, nil
}

func (c *MongoClient) DeleteStudio(id string) error {
	_, err := c.collectionStudio().DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		log.Error(err)
	}
	return nil
}

// studio에 player를 추가한다.이미 추가되어 있다면 update를 수행한다.
func (c *MongoClient) JoinPlayer(studioID string, playerID string) error {
	session, err := c.cli.StartSession()
	if err != nil {
		log.Error("StartSession:", err)
		return err
	}
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		filter := bson.M{"id": studioID}

		// delete old player info
		update := bson.M{
			"$pull": bson.M{"players": bson.M{"id": playerID}},
		}
		result := c.collectionStudio().FindOneAndUpdate(context.TODO(), filter, update)
		if result.Err() != nil && result.Err() != mongo.ErrNoDocuments {
			log.Error("detele old player:", result.Err())
			return nil, result.Err()
		}

		// insert new player info
		player := models.Player{
			ID:   playerID,
			Name: "", // TODO: 유저 DB에서 name을 가져와야 함.

		}
		update = bson.M{
			"$push": bson.M{"players": player},
		}
		result = c.collectionStudio().FindOneAndUpdate(context.TODO(), filter, update)
		if result.Err() != nil {
			log.Error("insert new player:", result.Err())
			return nil, result.Err()
		}

		return result, err
	}

	_, err = session.WithTransaction(context.TODO(), callback)
	if err != nil {
		log.Error("WithTransaction:", err)
		return err
	}

	return nil
}

func (c *MongoClient) LeavePlayer(studioID string, playerID string) error {
	filter := bson.M{"id": studioID}
	update := bson.M{
		"$pull": bson.M{"players": bson.M{"id": playerID}},
	}

	// TODO: not found id 처리
	result := c.collectionStudio().FindOneAndUpdate(context.TODO(), filter, update)
	if result.Err() != nil {
		log.Error("FindOneAndUpdate:", result.Err())
		return result.Err()
	}
	return nil
}

func (c *MongoClient) ListPlayers(studioID string) ([]models.Player, error) {
	type Result struct {
		Players []models.Player `bson:"players"`
	}
	var result Result

	filter := bson.M{"id": studioID}
	projection := bson.M{"players": 1}
	err := c.collectionStudio().FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		log.Error(err)
		return result.Players, err
	}

	return result.Players, nil
}

func (c *MongoClient) GetMQTT() (models.MQTTServer, error) {
	var mqtt models.MQTTServer
	err := c.collectionMqtt().FindOne(context.TODO(), bson.M{}).Decode(&mqtt)
	if err != nil {
		log.Error(err)
		return mqtt, err
	}

	return mqtt, nil
}

func (c *MongoClient) GetTURNs() ([]models.TURNServer, error) {
	var turns []models.TURNServer
	cursor, err := c.collectionTurn().Find(context.TODO(), bson.M{})
	if err != nil {
		log.Error(err)
		return turns, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &turns); err != nil {
		log.Error(err)
		return turns, err
	}

	return turns, nil
}
