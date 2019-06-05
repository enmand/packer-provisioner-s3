package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	awscommon "github.com/hashicorp/packer/builder/amazon/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/jbowes/vice"
)

const defaultTimeout = time.Minute * 2 // 2 minute timeout

// Config is the S3 provisioner configuration
type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	awscommon.AccessConfig `mapstructure:",squash"`

	Bucket    string `mapstructure:"bucket"`
	FileKey   string `mapstructure:"key"`
	LocalPath string `mapstructure:"local_path"`
	Timeout   string `mapstructure:"timeout"`

	timeout time.Duration

	ctx interpolate.Context
}

// Provisioner is a `file`-like Provisioner that can download a file on S3 to a
// remote instance
type Provisioner struct {
	packer.Provisioner

	config Config
	s3     *s3.S3
}

// Prepare prepres
func (p *Provisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return vice.Wrap(err, vice.InvalidArgument, "unable to decode options")
	}

	if p.config.Bucket == "" {
		return vice.Wrap(err, vice.InvalidArgument, "`bucket` is required")
	}

	if p.config.FileKey == "" {
		return vice.Wrap(err, vice.InvalidArgument, "`key` is required")
	}

	if p.config.LocalPath == "" {
		return vice.Wrap(err, vice.InvalidArgument, "`local_path` is required")
	}

	p.config.timeout, _ = time.ParseDuration(p.config.Timeout)
	if p.config.Timeout == "" {
		p.config.timeout = defaultTimeout
	}

	packer.LogSecretFilter.Set(p.config.AccessKey, p.config.SecretKey, p.config.Token)

	return nil
}

func (p *Provisioner) Provision(ctx context.Context, ui packer.Ui, comm packer.Communicator) error {
	ui.Say(fmt.Sprintf("Provisioning from S3..."))
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, p.config.timeout)
	defer cancel()

	session, err := p.config.Session()
	if err != nil {
		return vice.Wrap(err, vice.AuthRequired, "no aws session")
	}

	s3conn := s3.New(session)
	resp, err := s3conn.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.config.Bucket),
		Key:    aws.String(p.config.FileKey),
	})
	if err != nil {
		return vice.Wrap(err, vice.Temporary, "unable to download object from s3")
	}

	if err := comm.Upload(p.config.LocalPath, resp.Body, nil); err != nil {
		return vice.Wrap(err, vice.Temporary, "unable to upload file")
	}

	return nil
}

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterProvisioner(new(Provisioner))
}
