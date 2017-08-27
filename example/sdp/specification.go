package sdpExample

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
)

type testSpecification struct {
}

func (s *testSpecification) Concourse() (*primitive.Concourse, error) {
	return &primitive.Concourse{
		URL:      "http://concourse.com",
		User:     "user",
		Password: "password",
	}, nil
}

func (s *testSpecification) DeployImageRegistry() (*library.ImageRegistry, error) {
	return &library.ImageRegistry{
		Domain:             "registry.com",
		AwsAccessKeyId:     "key",
		AwsSecretAccessKey: "secret",
	}, nil
}

func (s *testSpecification) ConcourseBuilderGitSource() (*library.GitSource, error) {
	return &library.GitSource{
		Repo: &library.GitRepo{
			URI:        "git@github.com:concourse-friends/concourse-builder.git",
			PrivateKey: "private-key",
		},
		Branch: "master",
	}, nil
}

func (s *testSpecification) TargetGitRepo() (*library.GitRepo, error) {
	return &library.GitRepo{
		URI:        "git@github.com:target.git",
		PrivateKey: "private-key",
	}, nil
}

func (s *testSpecification) GenerateProjectLocation(resourceRegistry *project.ResourceRegistry, overrideBranch *primitive.GitBranch) (project.IRun, error) {
	return &primitive.Location{
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
