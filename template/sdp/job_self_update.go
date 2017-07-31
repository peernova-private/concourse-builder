package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

func SelfUpdateJob(pipelineLocation project.IRun) *project.Job {
	goImageResource := &project.JobResource{
		Name:    library.GoImage.Name,
		Trigger: true,
	}

	preparedDir := &library.TaskOutput{
		Directory: "prepared",
	}

	taskPrepare := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "prepare",
		Image:    goImageResource,
		Run:      pipelineLocation,
		Outputs: []project.IOutput{
			preparedDir,
		},
	}

	updateJob := &project.Job{
		Name:   project.JobName("self-update"),
		Groups: project.JobGroups{},
		Steps: project.ISteps{
			taskPrepare,
		},
	}
	return updateJob
}
