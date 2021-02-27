package routers

import (
	"api-server/controllers"
	"api-server/util"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/login", controllers.Login)
	app.Use(util.VerficateToken())
}

func SessionRoutes(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Post("/sessions", controllers.NewSession)
	v1.Get("/sessions", controllers.ListSessions)
	v1.Delete("/sessions/:id", controllers.DeleteSession)
}
