package resource

import (
	"github.com/concourse-friends/concourse-builder/model"
)

// Git resource source
type GitSource struct {
	// URI to the git repo
	URI string `yaml:",omitempty"`

	// branch name
	Branch string `yaml:",omitempty"`

	// Private key the has access to the repo
	PrivateKey string `yaml:"private_key,omitempty"`

	// Which paths to be checked for update
	Paths []string `yaml:",omitempty"`

	// Which paths to be excluded from being checked
	IgnorePaths []string `yaml:"ignore_paths,omitempty"`
}

// Git resource type
var GitResourceType = &model.ResourceType{
	// The name
	Name: "git",

	// The type
	Type: SystemResourceTypeName,
}

func init() {
	GlobalTypeRegistry.MustRegisterType(GitResourceType)
}
