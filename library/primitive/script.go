package primitive

import "github.com/concourse-friends/concourse-builder/project"

type Script struct {
	Location
	Dependencies []IVolume
}

func (s *Script) InputResources() project.JobResources {
	resources := s.Location.InputResources()
	for _, dependency := range s.Dependencies {
		if res, ok := dependency.(*project.JobResource); ok {
			resources = append(resources, res)
		}
	}
	return resources
}
