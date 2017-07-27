package library

import (
	"path"

	"github.com/concourse-friends/concourse-builder/project"
)

type FromParam project.Resource

func (fp *FromParam) Value() interface{} {
	source := fp.Source.(*ImageSource)
	result := path.Join(source.Repository.Domain, source.Location)
	if source.Tag != "" {
		result = result + ":" + source.Tag
	}
	return result
}
