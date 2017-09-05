package image

import (
	"fmt"
	"path"
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type Source struct {
	Registry   *Registry
	Repository string
	Tag        Tag
}

func (im *Source) ModelSource() interface{} {
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

