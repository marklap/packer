package arifactoverride

import (
	"fmt"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

// The artifact-override post-processor allows you to specify arbitrary files as
// artifacts. These will override any other artifacts created by the builder.
// This allows you to use a builder and provisioner to create some file, such as
// a compiled binary or tarball, extract it from the builder (VM or container)
// and then save that binary or tarball and throw away the builder.

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	Files []string `mapstructure:"files"`
	Keep  bool     `mapstructure:"keep_input_artifact"`

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	if len(p.config.Files) == 0 {
		return fmt.Errorf("No files specified in artifact-override; your build will not produce any artifacts")
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	artifact, err := NewArtifact(p.config.Files)

	return artifact, p.config.Keep, err
}
