package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type GitSource resource.GitSource

func (gs *GitSource) ModelSource() interface{} {
	return (*resource.GitSource)(gs)
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
