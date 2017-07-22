package sdpExample

import "github.com/concourse-friends/concourse-builder/resource"

type specification struct {
	ImageRepository *resource.ImageRepository
}

func (s *specification) DeployImageRepository() *resource.ImageRepository {
	return s.ImageRepository
}
