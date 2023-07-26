package routes

import (
	"github.com/gofiber/fiber/v2"
	"trendyol.com/security/appsec/devsecops/wolfichef/api/controllers"
)

func GitlabRoutes(router fiber.Router) {
	gitlab := router.Group("gitlab")
	gitlab.Get("url", controllers.GitlabOAuthUrl)
	gitlab.Post("token", controllers.GitlabFetchToken)
	gitlab.Post("refresh", controllers.GitlabRefreshToken)
}
