package controllers

import (
	"api-server/mongo"
	"api-server/util"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Unauthorize error
	mongoClient := mongo.NewMongoConn()
	if err := mongoClient.CheckAuth(user, pass); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	encodedToken, err := util.GenerateToken()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": encodedToken})
}
