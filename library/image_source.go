package library

import (
	"path"
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type ImageSource struct {
	Repository *ImageRepository
	Location   string
	Tag        string
}

func (im *ImageSource) ModelSource() interface{} {
	repository := im.Location
	if im.Repository.Domain != "" {
		repository = path.Join(im.Repository.Domain, repository)
	}

	return &resource.ImageSource{
		Repository: repository,
		Tag:        im.Tag,
	}
}

var UbuntuImage = &project.Resource{
	Name: "ubuntu-image",
	Type: resource.ImageResourceType.Name,
	Source: &ImageSource{
		Repository: DockerHub,
		Location:   "ubuntu",
		Tag:        "16.04",
	},
	CheckEvery: model.Duration(24 * time.Hour),
}
