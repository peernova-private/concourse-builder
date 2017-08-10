package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

func SelfUpdateJob(environment map[string]interface{}, concourse *library.Concourse, pipelineLocation project.IRun, flyImage project.ResourceName) *project.Job {
	flyImageResource := &project.JobResource{
		Name:    flyImage,
		Trigger: true,
	}

	taskCheck := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "check",
		Image:    flyImageResource,
		Run: &library.Location{
			Volume: &project.JobResource{
				Name:    library.ConcourseBuilderGitName,
				Trigger: true,
			},
			RelativePath: "scripts/check_fly_version.sh",
		},
		Params: map[string]interface{}{
			"CONCOURSE_URL": concourse.URL,
		},
	}

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
		Params:   environment,
		Outputs: []project.IOutput{
			pipelinesDir,
		},
	}

	taskPrepare.Params["PIPELINES"] = "pipelines"

	taskUpdate := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "update pipelines",
		Image:    flyImageResource,
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
			"CONCOURSE_URL":      concourse.URL,
			"CONCOURSE_USER":     concourse.User,
			"CONCOURSE_PASSWORD": concourse.Password,
		},
	}

	updateJob := &project.Job{
		Name:   project.JobName("self-update"),
		Groups: project.JobGroups{},
		Steps: project.ISteps{
			taskCheck,
			taskPrepare,
			taskUpdate,
		},
	}
	return updateJob
}
