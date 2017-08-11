package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type SdpSpecification interface {
	Concourse() (*library.Concourse, error)
	DeployImageRepository() (*library.ImageRegistry, error)
	ConcourseBuilderGitSource() (*library.GitSource, error)
	GenerateMainPipelineLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error)
	Environment() (map[string]interface{}, error)
}

func GenerateProject(specification SdpSpecification) (*project.Project, error) {
	mainPipeline := project.NewPipeline()
	mainPipeline.Name = "sdp"
	mainPipeline.AllJobsGroup = project.AllJobsGroupLast

	mainPipeline.ResourceRegistry.MustRegister(library.GoImage)
	mainPipeline.ResourceRegistry.MustRegister(library.UbuntuImage)

	concourseBuilderGitSource, err := specification.ConcourseBuilderGitSource()
	if err != nil {
		return nil, err
	}

	concourseBuilderGit := &project.Resource{
		Name:   library.ConcourseBuilderGitName,
		Type:   resource.GitResourceType.Name,
		Source: concourseBuilderGitSource,
	}

	mainPipeline.ResourceRegistry.MustRegister(concourseBuilderGit)

	dockerRepository, err := specification.DeployImageRepository()
	if err != nil {
		return nil, err
	}

	// Prepare curl image job
	curlImage := &project.Resource{
		Name: "curl-image",
		Type: resource.ImageResourceType.Name,
		Source: &library.ImageSource{
			Registry: dockerRepository,
			Location: "concourse-builder/curl-image",
		},
	}
	mainPipeline.ResourceRegistry.MustRegister(curlImage)

	curlImageJob := CurlImageJob(curlImage.Name)

	// Prepare fly image job
	flyImage := &project.Resource{
		Name: "fly-image",
		Type: resource.ImageResourceType.Name,
		Source: &library.ImageSource{
			Registry: dockerRepository,
			Location: "concourse-builder/fly-image",
		},
	}
	mainPipeline.ResourceRegistry.MustRegister(flyImage)

	concourse, err := specification.Concourse()
	if err != nil {
		return nil, err
	}

	flyImageJob := FlyImageJob(concourse, curlImage, flyImage.Name)
	flyImageJob.AddJobToRunAfter(curlImageJob)

	// Prepare git image job
	gitImage := &project.Resource{
		Name: "git-image",
		Type: resource.ImageResourceType.Name,
		Source: &library.ImageSource{
			Registry: dockerRepository,
			Location: "concourse-builder/git-image",
		},
	}
	mainPipeline.ResourceRegistry.MustRegister(gitImage)

	gitImageJob := GitImageJob(gitImage.Name)

	// Prepare self update job
	generateMainPipelineLocation, err := specification.GenerateMainPipelineLocation(mainPipeline.ResourceRegistry)
	if err != nil {
		return nil, err
	}

	environment, err := specification.Environment()
	if err != nil {
		return nil, err
	}

	selfUpdateJob := SelfUpdateJob(environment, concourse, generateMainPipelineLocation, flyImage.Name)
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
