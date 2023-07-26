package routes

import (
	"github.com/gofiber/fiber/v2"
	"trendyol.com/security/appsec/devsecops/wolfichef/api/controllers"
)

func WolfiRoutes(router fiber.Router) {
	wolfi := router.Group("wolfi")
	wolfi.Get("packages", controllers.WolfiPackages)
	wolfi.Post("build", controllers.WolfiBuild)
}
