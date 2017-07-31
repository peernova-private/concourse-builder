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
