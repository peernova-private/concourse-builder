package sdpBranchExample

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type testSpecification struct {
}

func (s *testSpecification) LinuxImage(resourceRegistry *project.ResourceRegistry) (*project.Resource, error) {
	return image.Ubuntu, nil
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

func (s *testSpecification) GoImage(resourceRegistry *project.ResourceRegistry) (*project.Resource, error) {
	return image.Go, nil
}

func (s *testSpecification) ConcourseBuilderGit() (*project.Resource, error) {
	return &project.Resource{
		Name: library.ConcourseBuilderGitName,
		Type: resource.GitResourceType.Name,
		Source: &library.GitSource{
			Repo: &primitive.GitRepo{
				URI:        "git@github.com:concourse-friends/concourse-builder.git",
				PrivateKey: "private-key",
			},
			Branch: &primitive.GitBranch{
				Name: "master",
			},
		},
	}, nil

}

func (s *testSpecification) GenerateProjectLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error) {
	gitResource, err := s.ConcourseBuilderGit()
	if err != nil {
		return nil, err
	}

	return &primitive.Location{
		Volume:       resourceRegistry.JobResource(gitResource, true, nil),
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
