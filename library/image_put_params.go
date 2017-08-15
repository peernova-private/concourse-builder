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
	Load      bool
	FromImage *project.JobResource
}

func (ipp *ImagePutParams) ModelParams() interface{} {
	params := &resource.ImagePutParams{
		Build:     ipp.Build.Path(),
		BuildArgs: ipp.BuildArgs,
	}

	if ipp.Load {
		params.LoadBase = ipp.FromImage.Path()
	}

	return params
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
