package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/trivy"

	"github.com/asaskevich/govalidator"
	"trendyol.com/security/appsec/devsecops/wolfichef/api/models"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/builder"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/config"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/config/models/parameters"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/index"
)

func WolfiPackages(c *fiber.Ctx) error {
	path, err := filepath.Abs("packages.json")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var packages index.Packages
	err = json.Unmarshal(file, &packages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	filteredPackages := make(index.Packages, 0)

	query := c.Query("query")

	for _, pkg := range packages {
		if strings.Contains(pkg.Name, query) {
			filteredPackages = append(filteredPackages, pkg)
		}
	}

	return c.JSON(filteredPackages)
}

func WolfiBuild(c *fiber.Ctx) error {
	var body models.Builder
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	_, errors := govalidator.ValidateStruct(body)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	var params []parameters.ParameterOption
	params = append(params, parameters.WithPackages(body.Packages))
	if body.WorkDir != "" {
		params = append(params, parameters.WithWorkDir(body.WorkDir))
	}
	if body.Entrypoint != "" {
		params = append(params, parameters.WithEntrypoint(body.Entrypoint))
	}
	if len(body.Environments) > 0 {
		env := make(map[string]string)
		for _, environment := range body.Environments {
			env[environment.Key] = environment.Value
		}
		params = append(params, parameters.WithEnvironment(env))
	}

	opts := config.BuildOptions(params...)
	cfg := config.Generate(opts)

	ctx, err := builder.NewContext(cfg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"step":    "context",
		})
	}
	defer os.RemoveAll(ctx.Options.WorkDir)

	layer, err := ctx.Build()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"step":    "build",
		})
	}

	file := ctx.CreateTempFile("tar.gz")
	err = ctx.BuildImage(fmt.Sprintf("%s:%s", body.ImageName, body.ImageTag), file, cfg, layer)
	image := trivy.Image{Name: file}

	err = image.Scan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"step":    "scan",
		})
	}

	result, err := image.Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"step":    "scan_result",
		})
	}

	critical, high := result.Categorize()
	if critical > 0 || high > 0 {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{
			"message": err.Error(),
			"step":    "vulnerabilities_found",
		})
	}

	digest, err := ctx.Publish(cfg, layer, body.ImageRepository, body.ImageName, body.ImageTag, body.Username, body.Token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"step":    "publish",
		})
	}

	return c.JSON(digest.Name())
}
