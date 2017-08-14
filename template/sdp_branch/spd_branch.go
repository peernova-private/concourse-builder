package sdpBranch

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

type Specification interface {
	BootstarpSpecification
}

func GenerateProject(specification Specification) (*project.Project, error) {
	mainPipeline := project.NewPipeline()
	mainPipeline.AllJobsGroup = project.AllJobsGroupLast

	concourseBuilderGitSource, err := specification.ConcourseBuilderGitSource()
	if err != nil {
		return nil, err
	}

	mainPipeline.Name = project.ConvertToPipelineName(concourseBuilderGitSource.Branch + "-sdpb")

	imageRegistry, err := specification.DeployImageRegistry()
	if err != nil {
		return nil, err
	}

	concourse, err := specification.Concourse()
	if err != nil {
		return nil, err
	}

	generateProjectLocation, err := specification.GenerateProjectLocation(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	environment, err := specification.Environment()
	if err != nil {
		return nil, err
	}

	selfUpdateJob := library.SelfUpdateJob(&library.SelfUpdateJobArgs{
		FlyImageJobArgs: &library.FlyImageJobArgs{
			CurlImageJobArgs: &library.CurlImageJobArgs{
				ConcourseBuilderGitSource: concourseBuilderGitSource,
				ImageRegistry:             imageRegistry,
				ResourceRegistry:          mainPipeline.ResourceRegistry,
				Tag:                       library.ConvertToImageTag(concourseBuilderGitSource.Branch),
			},
			Concourse: concourse,
		},
		Environment:             environment,
		GenerateProjectLocation: generateProjectLocation,
	})

	mainPipeline.Jobs = project.Jobs{
		selfUpdateJob,
	}

	prj := &project.Project{
		Pipelines: project.Pipelines{
			mainPipeline,
		},
	}

	return prj, nil
}
