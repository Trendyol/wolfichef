package builder

import (
	"context"
	"fmt"
	"os"

	"chainguard.dev/apko/pkg/build"
	"chainguard.dev/apko/pkg/build/oci"
	"chainguard.dev/apko/pkg/build/types"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Builder struct {
	*build.Context
}

func NewContext(config types.ImageConfiguration) (*Builder, error) {
	tmp, err := os.MkdirTemp("", "wolfi-*")
	if err != nil {
		return nil, err
	}

	builder, err := build.New(tmp, build.WithImageConfiguration(config))
	if err != nil {
		return nil, err
	}
	return &Builder{builder}, nil
}

func (bc *Builder) Build() (v1.Layer, error) {
	if err := bc.Refresh(); err != nil {
		return nil, err

	}

	_, layer, err := bc.BuildLayer()
	if err != nil {
		return nil, err

	}
	return layer, nil
}

func (bc *Builder) BuildImage(imageName string, output string, config types.ImageConfiguration, layer v1.Layer) error {
	err := oci.BuildImageTarballFromLayer(imageName, layer, output, config, bc.Logger(), bc.Options)
	if err != nil {
		return nil

	}
	return nil
}

func (bc *Builder) CreateTempFile(ext string) string {
	f, err := os.CreateTemp("/tmp", fmt.Sprintf("wolfi-*.%s", ext))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return f.Name()
}

func (bc *Builder) Publish(config types.ImageConfiguration, layer v1.Layer, host string, repository string, tag string, username string, token string) (name.Digest, error) {
	ctx := context.TODO()

	auth := remote.WithAuth(
		&authn.Basic{
			Username: username,
			Password: token,
		},
	)

	imgDigest, _, err := oci.PublishImageFromLayer(ctx, layer, config, bc.Options.SourceDateEpoch, config.Archs[0], bc.Logger(), false, true, []string{fmt.Sprintf("%s/%s:%s", host, repository, tag)}, auth)
	if err != nil {
		return name.Digest{}, err
	}
	return imgDigest, nil
}
