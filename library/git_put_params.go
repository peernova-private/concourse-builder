package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type IRepository interface {
	Path() string
}

type GitPutParams struct {
	Repository IRepository
	Force      bool
}

func (gpp *GitPutParams) ModelParams() interface{} {
	params := &resource.GitPutParams{
		Repository: gpp.Repository.Path(),
		Force:      gpp.Force,
	}

	return params
}

func (gpp *GitPutParams) InputResources() (project.JobResources, error) {
	var resources project.JobResources

	if res, ok := gpp.Repository.(project.IInputResource); ok {
		resources = append(resources, res.InputResources()...)
	}

	return resources, nil
}
