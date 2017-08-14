package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

type BranchesJobArgs struct {
	*library.GitImageJobArgs
	Environment             map[string]interface{}
	GenerateProjectLocation project.IRun
}

func BranchesJob(args *BranchesJobArgs) *project.Job {

	gitImage, _ := library.GitImageJob(args.GitImageJobArgs)

	gitImageResource := &project.JobResource{
		Name:    gitImage.Name,
		Trigger: true,
	}

	pipelinesDir := &library.TaskOutput{
		Directory: "pipelines",
	}

	taskUpdate := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "prepare pipelines",
		Image:    gitImageResource,
		Run:      args.GenerateProjectLocation,
		Params:   args.Environment,
		Outputs: []project.IOutput{
			pipelinesDir,
		},
	}

	branchesJob := &project.Job{
		Name:   project.JobName("branches"),
		Groups: project.JobGroups{},
		Steps: project.ISteps{
			taskUpdate,
		},
	}

	return branchesJob
}
