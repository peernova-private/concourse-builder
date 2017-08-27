package sdpBranch

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

type Specification interface {
	BootstrapSpecification
	ModifyJobs(resourceRegistry *project.ResourceRegistry) (project.Jobs, error)
	VerifyJobs(resourceRegistry *project.ResourceRegistry) (project.Jobs, error)
}

func GenerateProject(specification Specification) (*project.Project, error) {
	mainPipeline := project.NewPipeline()
	mainPipeline.AllJobsGroup = project.AllJobsGroupFirst
	mainPipeline.Name = project.ConvertToPipelineName(specification.Branch().Name() + "-sdpb")

	concourseBuilderGitSource, err := specification.ConcourseBuilderGitSource()
	if err != nil {
		return nil, err
	}

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

	modifyJobs, err := specification.ModifyJobs(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	for _, job := range modifyJobs {
		job.AddJobToRunAfter(selfUpdateJob)
	}

	mainPipeline.Jobs = append(mainPipeline.Jobs, modifyJobs...)

	verifyJobs, err := specification.VerifyJobs(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	for _, job := range verifyJobs {
		job.AddJobToRunAfter(selfUpdateJob)
		job.AddJobToRunAfter(modifyJobs...)
	}

	mainPipeline.Jobs = append(mainPipeline.Jobs, verifyJobs...)

	prj := &project.Project{
		Pipelines: project.Pipelines{
			mainPipeline,
		},
	}

	return prj, nil
}
