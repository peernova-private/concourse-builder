package sdpExample

import (
	"github.com/concourse-friends/concourse-builder/resource"
)

type testSpecification struct {
}

func (s *testSpecification) DeployImageRepository() *resource.ImageRepository {
	return &resource.ImageRepository{
		Domain: "repository.com",
	}
}

func (s *testSpecification) GitPrivateKey() (string, error) {
	return "private-key", nil
}
