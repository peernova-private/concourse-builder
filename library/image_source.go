package library

import (
	"fmt"
	"path"
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type ImageSource struct {
	Registry   *ImageRegistry
	Repository string
	Tag        ImageTag
}

func (im *ImageSource) ModelSource() interface{} {
	repository := im.Repository
	if im.Registry.Domain != "" {
		repository = path.Join(im.Registry.Domain, repository)
	}

	source := &resource.ImageSource{
		Repository: repository,
		Tag:        string(im.Tag),
	}

	if im.Registry.AwsAccessKeyId != "" || im.Registry.AwsSecretAccessKey != "" {
		if im.Registry.AwsAccessKeyId == "" || im.Registry.AwsSecretAccessKey == "" {
			return fmt.Errorf(
				"For ImageRegistry AwsAccessKeyId and AwsSecretAccessKey make sense only as pair",
				im.Registry.Domain)
		}

		source.AwsAccessKeyID = im.Registry.AwsAccessKeyId
		source.AwsSecretAccessKey = im.Registry.AwsSecretAccessKey
	}

	return source
}

var GoImage = &project.Resource{
	Name: "go-image",
	Type: resource.ImageResourceType.Name,
	Source: &ImageSource{
		Registry:   DockerHub,
		Repository: "golang",
		Tag:        "1.8",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}

var UbuntuImage = &project.Resource{
	Name: "ubuntu-image",
	Type: resource.ImageResourceType.Name,
	Source: &ImageSource{
		Registry:   DockerHub,
		Repository: "ubuntu",
		Tag:        "16.04",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}

var GradleImage = &project.Resource{
	Name: "gradle-image",
	Type: resource.ImageResourceType.Name,
	Source: &ImageSource{
		Registry:   DockerHub,
		Repository: "gradle",
	},
	CheckInterval: model.Duration(24 * time.Hour),
}
