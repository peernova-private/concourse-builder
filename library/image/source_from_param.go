package image

import (
	"path"

	"github.com/concourse-friends/concourse-builder/project"
)

type FromParam project.Resource

func (fp *FromParam) Value() string {
	source := fp.Source.(*Source)
	result := path.Join(source.Registry.Domain, source.Repository)
	if source.Tag != "" {
		result = result + ":" + string(source.Tag)
	}
	return result
}
