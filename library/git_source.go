package library

import (
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type GitTagFilter string

type GitSource struct {
	// Git repo and credentials
	Repo *primitive.GitRepo

	// branch name
	Branch *primitive.GitBranch

	// Git tag to sync to. The tag should be pointing to change in the branch
	TagFilter GitTagFilter
}

func (gs *GitSource) ModelSource() interface{} {
	return &resource.GitSource{
		URI:        gs.Repo.URI,
		PrivateKey: gs.Repo.PrivateKey,
		Branch:     gs.Branch.CanonicalName(),
		TagFilter:  string(gs.TagFilter),
	}
}

var ConcourseBuilderGitName project.ResourceName = "concourse-builder-git"
