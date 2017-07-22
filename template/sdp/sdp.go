package sdp

import (
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type SdpSpecification interface {
	DeployImageRepository() *resource.ImageRepository
}

var imagesGroup = &project.JobGroup{
	Name: "images",
}

var concourseBuilderGirResource = &model.Resource{
	Name: "concourse-builder-git",
	Type: resource.GitResourceType.Name,
	Source: &resource.GitSource{
		URI:    "git@github.com:concourse-friends/concourse-builder.git",
		Branch: "master",
	},
}

func GenerateProject(specification SdpSpecification) *project.Project {
	gitImageJob := GitImageJob(specification)

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
