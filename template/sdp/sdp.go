package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

type SdpSpecification interface {
	DeployImageRepository() *library.ImageRepository
	GitPrivateKey() (string, error)
}

var imagesGroup = &project.JobGroup{
	Name: "images",
}

func GenerateProject(specification SdpSpecification) (*project.Project, error) {
	gitImageJob, err := GitImageJob(specification)
	if err != nil {
		return nil, err
	}

	mainPipeline := &project.Pipeline{
		Jobs: project.Jobs{
			gitImageJob,
		},
	}

	prj := &project.Project{
		Pipelines: project.Pipelines{
			mainPipeline,
		},
	}

	return prj, nil
}
