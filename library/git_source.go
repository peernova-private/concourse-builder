package library

import (
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type GitSource struct {
	// Git repo and credentials
	Repo *primitive.GitRepo

	// branch name
	Branch *primitive.GitBranch
}

func (gs *GitSource) ModelSource() interface{} {
	return &resource.GitSource{
		URI:        gs.Repo.URI,
		PrivateKey: gs.Repo.PrivateKey,
		Branch:     gs.Branch.CanonicalName(),
	}
}

var ConcourseBuilderGitName project.ResourceName = "concourse-builder-git"

func RegisterConcourseBuilderGit(resourceRegistry *project.ResourceRegistry, source *GitSource) {
	if resourceRegistry.GetResource(ConcourseBuilderGitName) != nil {
		return
	}

	concourseBuilderGit := &project.Resource{
		Name:   ConcourseBuilderGitName,
		Type:   resource.GitResourceType.Name,
		Source: source,
	}

	resourceRegistry.MustRegister(concourseBuilderGit)
}
