package library

import (
	"path"

	"github.com/concourse-friends/concourse-builder/resource"
)

type ImageSource struct {
	Repository *ImageRepository
	Location   string
}

func (im *ImageSource) ModelSource() interface{} {
	return &resource.ImageSource{
		Repository: path.Join(im.Repository.Domain, im.Location),
	}
}
