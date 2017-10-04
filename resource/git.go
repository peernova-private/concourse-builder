package resource

import (
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
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

	// Optional. If specified as (list of pairs name and value) it will configure git global options,
	// setting each name with each value.
	GitConfig map[string]interface{} `yaml:"git_config,omitempty"`

	// Optional. If specified, the resource will only detect commits that have a tag matching
	// the expression that have been made against the branch. Patterns are glob(7) compatible
	// (as in, bash compatible).
	TagFilter string `yaml:"tag_filter,omitempty"`
}

type GitPutParams struct {
	Repository string `yaml:",omitempty"`
	Force      bool   `yaml:",omitempty"`
}

// Git resource type
var GitResourceType = &project.ResourceType{
	// The name
	Name: "git",

	// The type
	Type: model.SystemResourceTypeName,
}

func init() {
	project.GlobalTypeRegistry.MustRegisterType(GitResourceType)
}
