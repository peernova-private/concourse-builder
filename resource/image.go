package resource

import (
	"github.com/concourse-friends/concourse-builder/model"
)

// Image resource type
var ImageResourceType = &model.ResourceType{
	// The name
	Name: "docker-image",

	// The type
	Type: SystemResourceTypeName,
}

// Image resource source
type ImageSource struct {
	// Image repository
	Repository string

	// Image tag
	Tag string `yaml:",omitempty"`

	// Is the repo insecure
	Insecure bool `yaml:",omitempty"`

	// Optional. AWS access key to use for acquiring ECR credentials.
	AwsAccessKeyID string `yaml:"aws_access_key_id,omitempty"`

	// Optional. AWS secret key to use for acquiring ECR credentials.
	AwsSecretAccessKey string `yaml:"aws_secret_access_key,omitempty"`
}

type ImageGetParams struct {
	SkipDownload bool `yaml:"skip_download,omitempty"`
}

type ImagePutParams struct {
	Build     string
	BuildArgs map[string]interface{} `yaml:"build_args,omitempty"`
}

func init() {
	GlobalTypeRegistry.MustRegisterType(ImageResourceType)
}
