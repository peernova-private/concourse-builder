package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type SdpSpecification interface {
	DeployImageRepository() *library.ImageRepository
	GitPrivateKey() (string, error)
}

func GenerateProject(specification SdpSpecification) (*project.Project, error) {
	mainPipeline := project.NewPipeline()

	privateKey, err := specification.GitPrivateKey()
	if err != nil {
		return nil, err
	}

	concourseBuilderGit := &project.Resource{
		Name: library.ConcourseBuilderGit,
		Type: resource.GitResourceType.Name,
		Source: &library.GitSource{
			URI:        "git@github.com:concourse-friends/concourse-builder.git",
			Branch:     "master",
			PrivateKey: privateKey,
		},
	}

	mainPipeline.ResourceRegistry.MustRegister(concourseBuilderGit)

	gitImage := &project.Resource{
		Name: "git-image",
		Type: resource.ImageResourceType.Name,
		Source: &library.ImageSource{
			Repository: specification.DeployImageRepository(),
			Location:   "concourse-builder/git-image",
		},
	}
	mainPipeline.ResourceRegistry.MustRegister(library.UbuntuImage)
	mainPipeline.ResourceRegistry.MustRegister(gitImage)

	gitImageJob := GitImageJob(concourseBuilderGit.Name, gitImage.Name)

	mainPipeline.Jobs = project.Jobs{
		gitImageJob,
	}

	prj := &project.Project{
		Pipelines: project.Pipelines{
			mainPipeline,
		},
	}

	return prj, nil
}
