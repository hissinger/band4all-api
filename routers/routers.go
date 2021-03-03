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

func StudioRoutes(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Post("/studios", controllers.NewStudio)
	v1.Get("/studios", controllers.ListStudios)
	v1.Delete("/studios/:id", controllers.DeleteStudio)

	v1.Get("/studios/:sid/players", controllers.ListPlayers)
	v1.Put("/studios/:sid/players", controllers.JoinPlayer)
	v1.Delete("/studios/:sid/players/:pid", controllers.LeavePlayer)
}
