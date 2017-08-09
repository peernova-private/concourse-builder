package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

func SelfUpdateJob(privateKey string, pipelineLocation project.IRun, flyImage project.ResourceName) *project.Job {
	goImageResource := &project.JobResource{
		Name:    library.GoImage.Name,
		Trigger: true,
	}

	pipelinesDir := &library.TaskOutput{
		Directory: "pipelines",
	}

	taskPrepare := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "prepare pipelines",
		Image:    goImageResource,
		Run:      pipelineLocation,
		Params: map[string]interface{}{
			"CONCOURCE_BUILDER_GIT_PRIVATE_KEY": privateKey,
			"PIPELINES":                         "pipelines",
		},
		Outputs: []project.IOutput{
			pipelinesDir,
		},
	}

	taskUpdate := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "update pipelines",
		Image: &project.JobResource{
			Name:    flyImage,
			Trigger: true,
		},
		Run: &library.Location{
			Volume: &project.JobResource{
				Name:    library.ConcourseBuilderGitName,
				Trigger: true,
			},
			RelativePath: "scripts/set_pipelines.sh",
		},
		Params: map[string]interface{}{
			"PIPELINES": &library.Location{
				Volume: pipelinesDir,
			},
		},
	}

	updateJob := &project.Job{
		Name:   project.JobName("self-update"),
		Groups: project.JobGroups{},
		Steps: project.ISteps{
			taskPrepare,
			taskUpdate,
		},
	}
	return updateJob
}
