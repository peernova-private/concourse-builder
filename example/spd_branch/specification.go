package sdpBranchExample

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
)

type testSpecification struct {
}

func (s *testSpecification) Branch() *primitive.GitBranch {
	return &primitive.GitBranch{
		Name: "target_branch",
	}
}

func (s *testSpecification) Concourse() (*primitive.Concourse, error) {
	return &primitive.Concourse{
		URL:      "http://concourse.com",
		User:     "user",
		Password: "password",
	}, nil
}

func (s *testSpecification) DeployImageRegistry() (*image.Registry, error) {
	return &image.Registry{
		Domain:             "registry.com",
		AwsAccessKeyId:     "key",
		AwsSecretAccessKey: "secret",
	}, nil
}

func (s *testSpecification) GoImage() *project.Resource {
	return image.Go
}

func (s *testSpecification) ConcourseBuilderGitSource() (*library.GitSource, error) {
	return &library.GitSource{
		Repo: &primitive.GitRepo{
			URI:        "git@github.com:concourse-friends/concourse-builder.git",
			PrivateKey: "private-key",
		},
		Branch: &primitive.GitBranch{
			Name: "master",
		},
	}, nil
}

func (s *testSpecification) GenerateProjectLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error) {
	gitSource, err := s.ConcourseBuilderGitSource()
	if err != nil {
		return nil, err
	}

	library.RegisterConcourseBuilderGit(resourceRegistry, gitSource)

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

func (s *testSpecification) ModifyJobs(resourceRegistry *project.ResourceRegistry) (project.Jobs, error) {
	return nil, nil
}

func (s *testSpecification) VerifyJobs(resourceRegistry *project.ResourceRegistry) (project.Jobs, error) {
	return nil, nil
}
