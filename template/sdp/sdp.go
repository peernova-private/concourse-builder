package sdp

import "github.com/concourse-friends/concourse-builder/project"

type SdpSpecification interface {
}

var imagesGroup = &project.JobGroup{
	Name: "images",
}

func GenerateProject(specification SdpSpecification) *project.Project {

	gitImageJob := GitImageJob()

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

	return prj
}
