package routes

import "github.com/gofiber/fiber/v2"

func LoadRoutes(app *fiber.App) {
	api := app.Group("api")
	WolfiRoutes(api)
	GitlabRoutes(api)
}
