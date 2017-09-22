package resource

import (
	"github.com/concourse-friends/concourse-builder/project"
)

// Git pull request resource source
type GitPullRequestSource struct {
	// The repo name on github
	Repo string

	// URI to the git repo
	URI string `yaml:",omitempty"`

	// Private key that unlocks the repo
	PrivateKey string `yaml:"private_key,omitempty"`

	// Which is the branch that the pr is target to
	Base string `yaml:",omitempty"`

	// Github access token
	AccessToken string `yaml:"access_token,omitempty"`
}

// The pull request resource type
var PullRequestResourceType = &project.ResourceType{
	// The name
	Name: "pull-request",
	// The type
	Type: ImageResourceType.Type,

	// The official repo
	Source: &ImageSource{
		// Image repository
		Repository: "jtarchie/pr",
	},
}

func init() {
	project.GlobalTypeRegistry.MustRegisterType(PullRequestResourceType)
}
