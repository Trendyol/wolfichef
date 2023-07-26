package config

import (
	"chainguard.dev/apko/pkg/build/types"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/config/models/parameters"
)

func Generate(parameters parameters.Parameters) types.ImageConfiguration {
	config := types.ImageConfiguration{}

	config.Contents.Repositories = parameters.Repositories
	config.Contents.Packages = parameters.Packages
	config.Contents.Keyring = parameters.Keyring

	config.Accounts = parameters.Accounts
	config.Archs = parameters.Archs
	config.Environment = parameters.Environment

	if parameters.Entrypoint != "" {
		config.Entrypoint = types.ImageEntrypoint{Command: parameters.Entrypoint}
	}
	if parameters.WorkDir != "" {
		config.WorkDir = parameters.WorkDir
	}

	return config
}
