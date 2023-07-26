package main

import (
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"strings"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/builder"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/config"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/config/models/parameters"
)

type Arguments struct {
	Packages    *string `validate:"required"`
	Registry    *string `validate:"required,min=3"`
	Repository  *string `validate:"required,min=3"`
	Tag         *string `validate:"required,min=3"`
	DeployUser  *string `validate:"required,min=3"`
	DeployToken *string `validate:"required,min=3"`
}

func main() {
	var args Arguments
	args.Packages = flag.String("packages", "busybox", "comma seperated list of the packages you need")
	args.Registry = flag.String("registry", "localhost:5000", "registry url")
	args.Repository = flag.String("repository", "", "repository")
	args.Tag = flag.String("tag", "", "image tag")
	args.DeployUser = flag.String("user", "", "deploy user")
	args.DeployToken = flag.String("token", "", "deploy token")
	flag.Parse()

	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("%s field is required\n", err.StructField())
		}
		return
	}

	opts := config.BuildOptions(parameters.WithPackages(strings.Split(*args.Packages, ",")))
	cfg := config.Generate(opts)

	ctx, err := builder.NewContext(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(ctx.Options.WorkDir)

	layer, err := ctx.Build()
	if err != nil {
		log.Fatal(err)
	}

	digest, err := ctx.Publish(cfg, layer, *args.Registry, *args.Repository, *args.Tag, *args.DeployUser, *args.DeployToken)
	if err != nil {
		log.Fatal(err)
	}

	println(digest.Name())
}
