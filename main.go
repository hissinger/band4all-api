package main

import (
	"api-server/routers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n${body}\n${resBody}\n",
		Output: os.Stdout,
	}))

	routers.AuthRoutes(app)
	routers.StudioRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
