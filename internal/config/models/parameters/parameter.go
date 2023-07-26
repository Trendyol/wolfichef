package parameters

import (
	"chainguard.dev/apko/pkg/build/types"
)

type Parameters struct {
	Accounts     types.ImageAccounts
	Archs        []types.Architecture
	WorkDir      string
	Entrypoint   string
	Environment  map[string]string
	Keyring      []string
	Packages     []string
	Repositories []string
}

type ParameterOption func(*Parameters)

func WithRepositories(repositories []string) ParameterOption {
	return func(o *Parameters) {
		o.Repositories = repositories
	}
}

func WithPackages(packages []string) ParameterOption {
	return func(o *Parameters) {
		o.Packages = packages
	}
}

func WithEnvironment(environment map[string]string) ParameterOption {
	return func(o *Parameters) {
		o.Environment = environment
	}
}

func WithArchs(architectures []string) ParameterOption {
	var archs []types.Architecture
	return func(o *Parameters) {
		for arch := range architectures {
			archs = append(archs, types.Architecture(arch))
		}
		o.Archs = archs
	}
}

func WithAccounts(accounts types.ImageAccounts) ParameterOption {
	return func(o *Parameters) {
		o.Accounts = accounts
	}
}

func WithWorkDir(cwd string) ParameterOption {
	return func(o *Parameters) {
		o.WorkDir = cwd
	}
}

func WithEntrypoint(entrypoint string) ParameterOption {
	return func(o *Parameters) {
		o.Entrypoint = entrypoint
	}
}
