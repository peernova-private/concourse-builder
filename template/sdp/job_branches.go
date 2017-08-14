package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

type BranchesJobArgs struct {
	*library.GitImageJobArgs
	*library.FlyImageJobArgs
	Environment             map[string]interface{}
	GenerateProjectLocation project.IRun
}

func BranchesJob(args *BranchesJobArgs) *project.Job {

	gitImage, _ := library.GitImageJob(args.GitImageJobArgs)

	gitImageResource := &project.JobResource{
		Name:    gitImage.Name,
		Trigger: true,
	}

	branchesDir := &library.TaskOutput{
		Directory: "branches",
	}

	taskBranches := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "obtain branches",
		Image:    gitImageResource,
		Run: &library.Location{
			Volume: &library.Directory{
				Root: "/bin",
			},
			RelativePath: "obtain_branches.sh",
		},

		Params: args.Environment,
		Outputs: []project.IOutput{
			branchesDir,
		},
	}

	// TODO: Add task to create bootstrap pipelines for the branches

	flyImage, _ := library.FlyImageJob(args.FlyImageJobArgs)

	flyImageResource := &project.JobResource{
		Name:    flyImage.Name,
		Trigger: true,
	}

	pipelinesDir := &library.TaskOutput{
		Directory: "pipelines",
	}

	taskCreateMissing := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "create missing pipelines",
		Image:    flyImageResource,
		Run: &library.Location{
			Volume: &library.Directory{
				Root: "/bin",
			},
			RelativePath: "create_missing_pipelines.sh",
		},
		Params: map[string]interface{}{
			"PIPELINES": &library.Location{
				Volume: pipelinesDir,
			},
			"CONCOURSE_URL":      args.Concourse.URL,
			"CONCOURSE_USER":     args.Concourse.User,
			"CONCOURSE_PASSWORD": args.Concourse.Password,
		},
	}

	// TODO: Add task to remove pipeliens for already removed branches

	branchesJob := &project.Job{
		Name:   project.JobName("branches"),
		Groups: project.JobGroups{},
		Steps: project.ISteps{
			taskBranches,
			//taskCreatePipelines,
			taskCreateMissing,
			//taskRemoveNotNeeded,
		},
	}

	return branchesJob
}
