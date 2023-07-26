package main

import (
	"trendyol.com/security/appsec/devsecops/wolfichef/api"
	"trendyol.com/security/appsec/devsecops/wolfichef/api/config"
	"trendyol.com/security/appsec/devsecops/wolfichef/api/routes"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/index"
)

func main() {

	api.Http = &config.AppConfig{}
	api.Http.Setup()

	packages := index.Prepare()
	packages.Save()

	routes.LoadRoutes(api.Http.Server.App)
	api.Http.Server.Serve()

}
