package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type SdpSpecification interface {
	Concourse() (*library.Concourse, error)
	DeployImageRegistry() (*library.ImageRegistry, error)
	ConcourseBuilderGitSource() (*library.GitSource, error)
	GenerateMainPipelineLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error)
	Environment() (map[string]interface{}, error)
}

func GenerateProject(specification SdpSpecification) (*project.Project, error) {
	mainPipeline := project.NewPipeline()
	mainPipeline.AllJobsGroup = project.AllJobsGroupFirst

	mainPipeline.ResourceRegistry.MustRegister(library.GoImage)
	mainPipeline.ResourceRegistry.MustRegister(library.UbuntuImage)

	concourseBuilderGitSource, err := specification.ConcourseBuilderGitSource()
	if err != nil {
		return nil, err
	}

	mainPipeline.Name = project.ConvertToPipelineName(concourseBuilderGitSource.Branch + "-sdp")

	concourseBuilderGit := &project.Resource{
		Name:   library.ConcourseBuilderGitName,
		Type:   resource.GitResourceType.Name,
		Source: concourseBuilderGitSource,
	}

	mainPipeline.ResourceRegistry.MustRegister(concourseBuilderGit)

	imageRegistry, err := specification.DeployImageRegistry()
	if err != nil {
		return nil, err
	}

	// Prepare curl image job
	curlImage, curlImageJob := library.CurlImageJob(
		imageRegistry,
		mainPipeline.ResourceRegistry,
		library.ConvertToImageTag(concourseBuilderGitSource.Branch))

	// Prepare fly image job
	concourse, err := specification.Concourse()
	if err != nil {
		return nil, err
	}

	flyImage, flyImageJob := library.FlyImageJob(&library.FlyImageJobArgs{
		ImageRegistry:    imageRegistry,
		ResourceRegistry: mainPipeline.ResourceRegistry,
		Tag:              library.ConvertToImageTag(concourseBuilderGitSource.Branch),
		Concourse:        concourse,
		CurlResource:     curlImage,
	})
	flyImageJob.AddJobToRunAfter(curlImageJob)

	// Prepare git image job

	_, gitImageJob := library.GitImageJob(
		imageRegistry,
		mainPipeline.ResourceRegistry,
		library.ConvertToImageTag(concourseBuilderGitSource.Branch))

	// Prepare self update job
	generateMainPipelineLocation, err := specification.GenerateMainPipelineLocation(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	environment, err := specification.Environment()
	if err != nil {
		return nil, err
	}

	selfUpdateJob := library.SelfUpdateJob(&library.SelfUpdateJobArgs{
		Environment:      environment,
		Concourse:        concourse,
		PipelineLocation: generateMainPipelineLocation,
		FlyImage:         flyImage.Name,
	})
	selfUpdateJob.AddJobToRunAfter(flyImageJob)

	mainPipeline.Jobs = project.Jobs{
		curlImageJob,
		flyImageJob,
		gitImageJob,
		selfUpdateJob,
	}

	prj := &project.Project{
		Pipelines: project.Pipelines{
			mainPipeline,
		},
	}

	return prj, nil
}
