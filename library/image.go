package library

import (
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func BuildImage() project.IStep {
	image := &model.Resource{
		Type: resource.ImageResourceType.Name,
	}
	resource.GlobalRegistry.MustRegisterResource(image)

	step := &project.PutStep{
		Resource: image,
	}

	return step
}
