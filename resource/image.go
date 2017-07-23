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
}

type ImagePutParams struct {
	Build string
}

func init() {
	GlobalTypeRegistry.MustRegisterType(ImageResourceType)
}
