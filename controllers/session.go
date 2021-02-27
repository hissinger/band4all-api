package controllers

import (
	"api-server/models"
	"api-server/mongo"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type session struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Creator     string `json:"creator"`
}

func NewSession(c *fiber.Ctx) error {
	// New Employee struct
	s := &session{}

	// Parse body into struct
	if err := c.BodyParser(s); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	mongoClient := mongo.NewMongoConn()
	mqtt, err := mongoClient.GetMQTT()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	turns, err := mongoClient.GetTURNs()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	sess := models.Session{
		ID:          uuid.New().String(),
		Title:       s.Title,
		Description: s.Description,
		Private:     s.Private,
		Creator:     s.Creator,
		CreatedDate: time.Now(),
	}
	if err := mongoClient.CreateSession(sess); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": sess.ID, "mqtt": mqtt, "turn": turns})
}

func ListSessions(c *fiber.Ctx) error {
	mongoClient := mongo.NewMongoConn()
	sessions, err := mongoClient.ListSessions()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"sessions": sessions})
}

func DeleteSession(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	mongoClient := mongo.NewMongoConn()
	err := mongoClient.DeleteSession(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
