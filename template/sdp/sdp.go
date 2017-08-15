package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/template/sdp_branch"
)

type Specification interface {
	Concourse() (*library.Concourse, error)
	DeployImageRegistry() (*library.ImageRegistry, error)
	ConcourseBuilderGitSource() (*library.GitSource, error)
	GenerateProjectLocation(resourceRegistry *project.ResourceRegistry, overrideBranch string) (project.IRun, error)
	Environment() (map[string]interface{}, error)
	BootstrapBranches() []string
}

func GenerateProject(specification Specification) (*project.Project, error) {
	prj := &project.Project{}

	for _, branch := range specification.BootstrapBranches() {
		branchSpecification := &BranchBootstrapSpecification{
			Specification: specification,
			Branch:        branch,
		}
		project, err := sdpBranch.GenerateBootstarpProject(branchSpecification)
		if err != nil {
			return nil, err
		}
		prj.Pipelines = append(prj.Pipelines, project.Pipelines...)
	}

	mainPipeline := project.NewPipeline()
	mainPipeline.AllJobsGroup = project.AllJobsGroupFirst

	concourseBuilderGitSource, err := specification.ConcourseBuilderGitSource()
	if err != nil {
		return nil, err
	}

	mainPipeline.Name = project.ConvertToPipelineName(concourseBuilderGitSource.Branch + "-sdp")

	imageRegistry, err := specification.DeployImageRegistry()
	if err != nil {
		return nil, err
	}

	concourse, err := specification.Concourse()
	if err != nil {
		return nil, err
	}

	generateProjectLocation, err := specification.GenerateProjectLocation(mainPipeline.ResourceRegistry, "")
	if err != nil {
		return nil, err
	}

	environment, err := specification.Environment()
	if err != nil {
		return nil, err
	}

	curlImageJobArgs := &library.CurlImageJobArgs{
		ConcourseBuilderGitSource: concourseBuilderGitSource,
		ImageRegistry:             imageRegistry,
		ResourceRegistry:          mainPipeline.ResourceRegistry,
		Tag:                       library.ConvertToImageTag(concourseBuilderGitSource.Branch),
	}

	flyImageJobArgs := &library.FlyImageJobArgs{
		CurlImageJobArgs: curlImageJobArgs,
		Concourse:        concourse,
	}

	selfUpdateJob := library.SelfUpdateJob(&library.SelfUpdateJobArgs{
		FlyImageJobArgs:         flyImageJobArgs,
		Environment:             environment,
		GenerateProjectLocation: generateProjectLocation,
	})

	branchesJob := BranchesJob(&BranchesJobArgs{
		GitImageJobArgs: &library.GitImageJobArgs{
			CurlImageJobArgs: curlImageJobArgs,
		},
		FlyImageJobArgs:         flyImageJobArgs,
		Environment:             environment,
		GenerateProjectLocation: generateProjectLocation,
	})

	mainPipeline.Jobs = project.Jobs{
		selfUpdateJob,
		branchesJob,
	}

	prj.Pipelines = append(prj.Pipelines, mainPipeline)
	return prj, nil
}
