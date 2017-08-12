package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type IBuild interface {
	Path() string
	InputResources() project.JobResources
}

type ImagePutParams struct {
	Build     IBuild
	BuildArgs map[string]interface{}
	FromImage *project.JobResource
}

func (ipp *ImagePutParams) ModelParams() interface{} {
	return &resource.ImagePutParams{
		Build:     ipp.Build.Path(),
		BuildArgs: ipp.BuildArgs,
	}
}

func (ipp *ImagePutParams) InputResources() (project.JobResources, error) {
	var resources project.JobResources

	res := ipp.Build.InputResources()
	resources = append(resources, res...)

	if ipp.FromImage != nil {
		resources = append(resources, ipp.FromImage)
	}

	return resources, nil
}
