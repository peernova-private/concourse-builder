package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

type BranchBootstrapSpecification struct {
	Specification Specification
	Branch        string
}

func (bbs *BranchBootstrapSpecification) Concourse() (*library.Concourse, error) {
	return bbs.Specification.Concourse()
}

func (bbs *BranchBootstrapSpecification) DeployImageRegistry() (*library.ImageRegistry, error) {
	return bbs.Specification.DeployImageRegistry()
}

func (bbs *BranchBootstrapSpecification) ConcourseBuilderGitSource() (*library.GitSource, error) {
	return bbs.Specification.ConcourseBuilderGitSource()
}

func (bbs *BranchBootstrapSpecification) GenerateProjectLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error) {
	return bbs.Specification.GenerateProjectLocation(resourceRegistry, bbs.Branch)
}

func (bbs *BranchBootstrapSpecification) Environment() (map[string]interface{}, error) {
	return bbs.Specification.Environment()
}
