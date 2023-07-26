package config

import (
	"chainguard.dev/apko/pkg/build/types"
	"trendyol.com/security/appsec/devsecops/wolfichef/internal/config/models/parameters"
)

func BuildOptions(options ...parameters.ParameterOption) parameters.Parameters {
	opts := parameters.Parameters{
		Accounts: types.ImageAccounts{
			Groups: []types.Group{
				{
					GroupName: "trend",
					GID:       65532,
				},
			},
			Users: []types.User{
				{
					UserName: "trend",
					UID:      65532,
					GID:      65532,
				},
			},
			RunAs: "65532",
		},
		Archs: []types.Architecture{
			types.Architecture("arm64"),
		},
		Entrypoint:   "/bin/sh",
		Environment:  map[string]string{"PATH": "/usr/sbin:/sbin:/usr/bin:/bin"},
		Keyring:      []string{"https://packages.wolfi.dev/bootstrap/stage3/wolfi-signing.rsa.pub"},
		Repositories: []string{"https://packages.wolfi.dev/os"},
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}
