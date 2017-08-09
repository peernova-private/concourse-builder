package library

import (
	"path"
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type ImageSource struct {
	Registry *ImageRegistry
	Location string
	Tag      string
}

func (im *ImageSource) ModelSource() interface{} {
	repository := im.Location
	if im.Registry.Domain != "" {
		repository = path.Join(im.Registry.Domain, repository)
	}

	return &resource.ImageSource{
		Repository: repository,
		Tag:        im.Tag,
	}
}

var GoImage = &project.Resource{
	Name: "go-image",
	Type: resource.ImageResourceType.Name,
	Source: &ImageSource{
		Registry: DockerHub,
		Location: "golang",
		Tag:      "1.8",
	},
	CheckEvery: model.Duration(24 * time.Hour),
}

var UbuntuImage = &project.Resource{
	Name: "ubuntu-image",
	Type: resource.ImageResourceType.Name,
	Source: &ImageSource{
		Registry: DockerHub,
		Location: "ubuntu",
		Tag:      "16.04",
	},
	CheckEvery: model.Duration(24 * time.Hour),
}
