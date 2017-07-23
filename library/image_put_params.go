package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type ImagePutParams struct {
	Build *Location
}

func (ipp *ImagePutParams) ModelParams() interface{} {
	return &resource.ImagePutParams{
		Build: ipp.Build.String(),
	}
}

func (ipp *ImagePutParams) InputResources() (project.JobResources, error) {
	var resources project.JobResources

	res, ok := ipp.Build.Volume.(*project.JobResource)
	if ok {
		resources = append(resources, res)
	}
	return resources, nil
}
