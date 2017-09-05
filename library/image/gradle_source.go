package image

import (
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

var Gradle = &project.Resource{
	Name: "gradle-image",
	Type: resource.ImageResourceType.Name,
	Source: &Source{
		Registry:   DockerHub,
		Repository: "gradle",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}
