package sdp

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type SdpSpecification interface {
	DeployImageRepository() *resource.ImageRepository
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
