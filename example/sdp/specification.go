package sdpExample

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

type testSpecification struct {
}

func (s *testSpecification) FlyVersion() (string, error) {
	return "v3.3.1", nil
}

func (s *testSpecification) DeployImageRepository() (*library.ImageRegistry, error) {
	return &library.ImageRegistry{
		Domain: "registry.com",
	}, nil
}

func (s *testSpecification) ConcourseBuilderGitPrivateKey() (string, error) {
	return "private-key", nil
}

func (s *testSpecification) GenerateMainPipelineLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error) {
	return &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "foo",
	}, nil
}
