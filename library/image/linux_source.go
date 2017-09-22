package image

import (
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

var Alpine = &project.Resource{
	Name: "alpine-image",
	Type: resource.ImageResourceType.Name,
	Source: &Source{
		Registry:   DockerHub,
		Repository: "alpine",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}

var Ubuntu = &project.Resource{
	Name: "ubuntu-image",
	Type: resource.ImageResourceType.Name,
	Source: &Source{
		Registry:   DockerHub,
		Repository: "ubuntu",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}

var Ubuntu1604 = &project.Resource{
	Name: "ubuntu-image",
	Type: resource.ImageResourceType.Name,
	Source: &Source{
		Registry:   DockerHub,
		Repository: "ubuntu",
		Tag:        "16.04",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}
