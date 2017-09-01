package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
)

type BranchBootstrapSpecification struct {
	Specification Specification
	TargetBranch  *primitive.GitBranch
}

func (bbs *BranchBootstrapSpecification) Branch() *primitive.GitBranch {
	return bbs.TargetBranch
}

func (bbs *BranchBootstrapSpecification) Concourse() (*primitive.Concourse, error) {
	return bbs.Specification.Concourse()
}

func (bbs *BranchBootstrapSpecification) DeployImageRegistry() (*image.Registry, error) {
	return bbs.Specification.DeployImageRegistry()
}

func (bbs *BranchBootstrapSpecification) GoImage() (*project.Resource, error) {
	return bbs.Specification.GoImage(), nil
}

func (bbs *BranchBootstrapSpecification) ConcourseBuilderGitSource() (*library.GitSource, error) {
	return bbs.Specification.ConcourseBuilderGitSource()
}

func (bbs *BranchBootstrapSpecification) GenerateProjectLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error) {
	return bbs.Specification.GenerateProjectLocation(resourceRegistry, bbs.TargetBranch)
}

func (bbs *BranchBootstrapSpecification) Environment() (map[string]interface{}, error) {
	enviroment, err := bbs.Specification.Environment()
	if err != nil {
		return nil, err
	}
	enviroment["BRANCH"] = bbs.TargetBranch.CanonicalName()
	return enviroment, nil
}
