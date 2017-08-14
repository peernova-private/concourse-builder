package sdpExample

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

type testSpecification struct {
}

func (s *testSpecification) Concourse() (*library.Concourse, error) {
	return &library.Concourse{
		URL:      "http://concourse.com",
		User:     "user",
		Password: "password",
	}, nil
}

func (s *testSpecification) DeployImageRegistry() (*library.ImageRegistry, error) {
	return &library.ImageRegistry{
		Domain: "registry.com",
	}, nil
}

func (s *testSpecification) ConcourseBuilderGitSource() (*library.GitSource, error) {
	return &library.GitSource{
		URI:        "git@github.com:concourse-friends/concourse-builder.git",
		Branch:     "master",
		PrivateKey: "private-key",
	}, nil
}

func (s *testSpecification) GenerateProjectLocation(resourceRegistry *project.ResourceRegistry, overrideBranch string) (project.IRun, error) {
	return &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "foo",
	}, nil
}

func (s *testSpecification) Environment() (map[string]interface{}, error) {
	return map[string]interface{}{
		"BRANCH": "branch",
	}, nil
}

func (s *testSpecification) BootstrapBranches() []string {
	return nil
}
