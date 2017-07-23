package sdpExample

import "github.com/concourse-friends/concourse-builder/library"

type testSpecification struct {
}

func (s *testSpecification) DeployImageRepository() *library.ImageRepository {
	return &library.ImageRepository{
		Domain: "repository.com",
	}
}

func (s *testSpecification) GitPrivateKey() (string, error) {
	return "private-key", nil
}
