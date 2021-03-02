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

type join struct {
	ID string `json:"id"`
}

func NewSession(c *fiber.Ctx) error {
	s := &session{}
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
		Members:     []string{},
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

func JoinSession(c *fiber.Ctx) error {
	sessionID := c.Params("sid")
	if sessionID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	body := &join{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	mongoClient := mongo.NewMongoConn()
	err := mongoClient.JoinSession(sessionID, body.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func LeaveSession(c *fiber.Ctx) error {
	sessionID := c.Params("sid")
	if sessionID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	memberID := c.Params("mid")
	if memberID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	mongoClient := mongo.NewMongoConn()
	err := mongoClient.LeaveSession(sessionID, memberID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
