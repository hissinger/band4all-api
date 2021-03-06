package mongo

import (
	"api-server/models"
	"context"
	"strings"

	"github.com/iegomez/mosquitto-go-auth/hashing"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mqtt acl를 추가한다. 이미 username의 acl이 존재하면, update를 수행한다.
func (c *MongoClient) AddUserAcl(studioID string, username string, password string, superuser bool) error {
	user := models.MqttUser{
		Username:  username,
		Password:  hashPassword(password),
		Superuser: superuser,
		Acls:      []models.MqttAcl{},
	}
	user.Acls = append(user.Acls, models.MqttAcl{Topic: strings.Join([]string{studioID, "#"}, "/"), Acc: 3})

	filter := bson.M{"username": username}
	opts := options.Replace().SetUpsert(true)
	_, err := c.cli.Database("mosquitto").Collection("users").ReplaceOne(context.TODO(), filter, user, opts)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error("ReplaceOne fail:", err)
		return err
	}

	return nil
}

func (c *MongoClient) DelUserAcl(username string) error {
	_, err := c.cli.Database("mosquitto").Collection("users").DeleteOne(context.TODO(), bson.M{"username": username})
	if err != nil {
		log.Error("DeleteOne() fail:", err)
		return err
	}

	return nil
}

func hashPassword(password string) string {
	saltSize := 16
	saltEncoding := "base64"
	algorithm := "sha512"
	shaSize := hashing.SHA512Size
	iterations := 100000

	hashComparer := hashing.NewPBKDF2Hasher(saltSize, iterations, algorithm, saltEncoding, shaSize)
	pwHash, err := hashComparer.Hash(password)
	if err != nil {
		log.Error("password hash fail:", err)
	}
	return pwHash
}
