package library

import (
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

type SelfUpdateJobArgs struct {
	*FlyImageJobArgs
	Environment             map[string]interface{}
	GenerateProjectLocation project.IRun
}

func SelfUpdateJob(args *SelfUpdateJobArgs) *project.Job {
	flyImage, _ := FlyImageJob(args.FlyImageJobArgs)

	flyImageResource := &project.JobResource{
		Name:    flyImage.Name,
		Trigger: true,
	}

	taskCheck := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "check",
		Image:    flyImageResource,
		Run: &Location{
			Volume: &Directory{
				Root: "/bin/fly",
			},
			RelativePath: "check_version.sh",
		},
		Environment: map[string]interface{}{},
	}

	args.Concourse.PublicAccessEnvironment(taskCheck.Environment)

	args.ResourceRegistry.MustRegister(GoImage)

	goImageResource := &project.JobResource{
		Name:    GoImage.Name,
		Trigger: true,
	}

	pipelinesDir := &TaskOutput{
		Directory: "pipelines",
	}

	taskPrepare := &project.TaskStep{
		Platform:    model.LinuxPlatform,
		Name:        "prepare pipelines",
		Image:       goImageResource,
		Run:         args.GenerateProjectLocation,
		Environment: args.Environment,
		Outputs: []project.IOutput{
			pipelinesDir,
		},
	}

	taskPrepare.Environment["PIPELINES"] = "pipelines"

	taskUpdate := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "update pipelines",
		Image:    flyImageResource,
		Run: &Location{
			Volume: &Directory{
				Root: "/bin/fly",
			},
			RelativePath: "set_pipelines.sh",
		},
		Environment: map[string]interface{}{
			"PIPELINES": &Location{
				Volume: pipelinesDir,
			},
		},
	}
	args.Concourse.Environment(taskUpdate.Environment)

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
