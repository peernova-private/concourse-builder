package image

import (
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

var Go = &project.Resource{
	Name: "go-image",
	Type: resource.ImageResourceType.Name,
	Source: &Source{
		Registry:   DockerHub,
		Repository: "golang",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}

var Go18x = &project.Resource{
	Name: "go-image",
	Type: resource.ImageResourceType.Name,
	Source: &Source{
		Registry:   DockerHub,
		Repository: "golang",
		Tag:        "1.8",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}

var Go183 = &project.Resource{
	Name: "go-image",
	Type: resource.ImageResourceType.Name,
	Source: &Source{
		Registry:   DockerHub,
		Repository: "golang",
		Tag:        "1.8.3",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}
