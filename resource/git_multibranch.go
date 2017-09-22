package resource

import (
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

type GitMultibranchSource struct {
	// URI to the git repo
	URI string `yaml:",omitempty"`

	// Private key the has access to the repo
	PrivateKey string `yaml:"private_key,omitempty"`

	// Optional. Turns on multi-branch mode; takes a regular expression as argument --
	// branches matching the regular expression on origin will all be checked for changes.
	// Uses grep-style regular expression syntax
	Branches string
}

// The git multibranch resource type
var GitMultibranchResourceType = &project.ResourceType{
	// The name
	Name: "git-multibranch",
	// The type
	Type: model.DockerImageType,

	// The official repo
	Source: &ImageSource{
		// Image repository
		Repository: "cfcommunity/git-multibranch-resource",
	},
}

func init() {
	project.GlobalTypeRegistry.MustRegisterType(GitMultibranchResourceType)
}
