package controllers

import (
	"api-server/models"
	"api-server/mongo"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type studio struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Creator     string `json:"creator"`
}

type join struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func NewStudio(c *fiber.Ctx) error {
	s := &studio{}
	if err := c.BodyParser(s); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	mongoClient := mongo.NewMongoConn()
	defer mongoClient.Close()

	mqtt, err := mongoClient.GetMQTT()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	turns, err := mongoClient.GetTURNs()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	studio := models.Studio{
		ID:          uuid.New().String(),
		Title:       s.Title,
		Description: s.Description,
		Private:     s.Private,
		Creator:     s.Creator,
		CreatedDate: time.Now(),
		Players:     []models.Player{},
	}
	if err := mongoClient.CreateStudio(studio); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": studio.ID, "mqtt": mqtt, "turn": turns})
}

func ListStudios(c *fiber.Ctx) error {
	mongoClient := mongo.NewMongoConn()
	defer mongoClient.Close()

	studios, err := mongoClient.ListStudios()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"studios": studios})
}

func DeleteStudio(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	mongoClient := mongo.NewMongoConn()
	defer mongoClient.Close()

	err := mongoClient.DeleteStudio(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func JoinPlayer(c *fiber.Ctx) error {
	studioID := c.Params("sid")
	if studioID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	body := &join{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	mongoClient := mongo.NewMongoConn()
	defer mongoClient.Close()

	if err := mongoClient.JoinPlayer(studioID, body.ID); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := mongoClient.AddUserAcl(studioID, body.ID, body.Password, false); err != nil {
		mongoClient.LeavePlayer(studioID, body.ID)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func ListPlayers(c *fiber.Ctx) error {
	studioID := c.Params("sid")
	if studioID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mongoClient := mongo.NewMongoConn()
	defer mongoClient.Close()

	players, err := mongoClient.ListPlayers(studioID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"players": players})
}

func LeavePlayer(c *fiber.Ctx) error {
	studioID := c.Params("sid")
	if studioID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	playerID := c.Params("pid")
	if playerID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	mongoClient := mongo.NewMongoConn()
	defer mongoClient.Close()

	if err := mongoClient.DelUserAcl(playerID); err != nil {
		log.Warning(err)
	}

	if err := mongoClient.LeavePlayer(studioID, playerID); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
