package sdp

import "github.com/ggeorgiev/concourse-builder/project"

type SdpSpecification interface {
}

func GenerateProject(specification SdpSpecification) *project.Project {
	imagesGroup := &project.JobGroup{
		Name: "images",
	}

	gitImageJob := &project.Job{
		Name: "git-image",
		Groups: project.JobGroups{
			imagesGroup,
		},
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

	return prj
}
