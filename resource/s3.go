package resource

import (
	"github.com/concourse-friends/concourse-builder/model"
)

// S3 resource source
type S3Source struct {
	// S3 bucket
	Bucket string `yaml:",omitempty"`

	// S3 acess key id
	AccessKeyID string `yaml:"access_key_id"`

	// S3 secret access key
	SecretAccessKey string `yaml:"secret_access_key"`

	// S3 region name
	RegionName string `yaml:"region_name"`

	// Matching regex
	Regex string `yaml:",omitempty"`

	// S3 versioned file
	VersionedFile string `yaml:"versioned_file,omitempty"`
}

type S3GetParams struct {
	Unpack bool `yaml:",omitempty"`
}

type S3PutParams struct {
	File string `yaml:",omitempty"`
}

// S3 resource type
var S3ResourceType = &model.ResourceType{
	// The name
	Name: "s3",

	// The type
	Type: SystemResourceTypeName,
}

func init() {
	GlobalTypeRegistry.MustRegisterType(S3ResourceType)
}
