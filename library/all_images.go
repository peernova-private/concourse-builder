package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/jinzhu/copier"
)

type AllImagesArgs struct {
	LinuxImageResource  *project.Resource
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
	Concourse           *primitive.Concourse
}

func AllImages(args *AllImagesArgs) *project.Job {
	flyImageJobArgs := &FlyImageJobArgs{}
	copier.Copy(flyImageJobArgs, args)
	flyImage := FlyImageJob(flyImageJobArgs)
	flyImageResource := args.ResourceRegistry.JobResource(flyImage, true, nil)

	gitImageJobArgs := &GitImageJobArgs{}
	copier.Copy(gitImageJobArgs, args)
	gitImage := GitImageJob(gitImageJobArgs)
	gitImageResource := args.ResourceRegistry.JobResource(gitImage, true, nil)

	taskDummy := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "dummy",
		Image:    flyImageResource,
		Run: &primitive.Location{
			Volume: &primitive.Directory{
				Root: "/bin",
			},
			RelativePath: "echo",
		},
		Environment: map[string]interface{}{
			"FLY_IMAGE": &primitive.Location{
				Volume: flyImageResource,
			},
			"GIT_IMAGE": &primitive.Location{
				Volume: gitImageResource,
			},
		},
	}

	dummyJob := &project.Job{
		Name:   project.JobName("dummy"),
		Groups: project.JobGroups{},
		Steps: project.ISteps{
			taskDummy,
		},
	}

	return dummyJob
}
