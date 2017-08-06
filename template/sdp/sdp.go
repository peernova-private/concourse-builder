package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type SdpSpecification interface {
	FlyVersion() (string, error)
	DeployImageRepository() (*library.ImageRepository, error)
	ConcourseBuilderGitPrivateKey() (string, error)
	GenerateMainPipelineLocation(resourceRegistry *project.ResourceRegistry) (project.IRun, error)
}

func GenerateProject(specification SdpSpecification) (*project.Project, error) {
	mainPipeline := project.NewPipeline()
	mainPipeline.Name = "sdp"
	mainPipeline.AllJobsGroup = project.AllJobsGroupLast

	mainPipeline.ResourceRegistry.MustRegister(library.GoImage)
	mainPipeline.ResourceRegistry.MustRegister(library.UbuntuImage)

	privateKey, err := specification.ConcourseBuilderGitPrivateKey()
	if err != nil {
		return nil, err
	}

	concourseBuilderGit := &project.Resource{
		Name: library.ConcourseBuilderGitName,
		Type: resource.GitResourceType.Name,
		Source: &library.GitSource{
			URI:        "git@github.com:concourse-friends/concourse-builder.git",
			Branch:     "master",
			PrivateKey: privateKey,
		},
	}

	mainPipeline.ResourceRegistry.MustRegister(concourseBuilderGit)

	dockerRepository, err := specification.DeployImageRepository()
	if err != nil {
		return nil, err
	}

	// Prepare fly image job
	flyImage := &project.Resource{
		Name: "fly-image",
		Type: resource.ImageResourceType.Name,
		Source: &library.ImageSource{
			Repository: dockerRepository,
			Location:   "concourse-builder/fly-image",
		},
	}
	mainPipeline.ResourceRegistry.MustRegister(flyImage)

	flyVersion, err := specification.FlyVersion()
	if err != nil {
		return nil, err
	}

	flyImageJob := FlyImageJob(flyVersion, flyImage.Name)

	// Prepare git image job
	gitImage := &project.Resource{
		Name: "git-image",
		Type: resource.ImageResourceType.Name,
		Source: &library.ImageSource{
			Repository: dockerRepository,
			Location:   "concourse-builder/git-image",
		},
	}
	mainPipeline.ResourceRegistry.MustRegister(gitImage)

	gitImageJob := GitImageJob(gitImage.Name)

	// Prepare self update job
	generateMainPipelineLocation, err := specification.GenerateMainPipelineLocation(mainPipeline.ResourceRegistry)

	selfUpdateJob := SelfUpdateJob(privateKey, generateMainPipelineLocation, flyImage.Name)
	selfUpdateJob.AddJobToRunAfter(flyImageJob)

	mainPipeline.Jobs = project.Jobs{
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
