package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

var ImagesGroup = &project.JobGroup{
	Name: "images",
	Before: project.JobGroups{
		project.SystemGroup,
	},
}

type BuildImageArgs struct {
	Prepare            *project.Resource
	From               *project.Resource
	Name               string
	DockerFileResource project.IEnvironmentValue
	Image              project.ResourceName
	BuildArgs          map[string]interface{}
	PreprepareSteps    project.ISteps
	Eval               string
}

func BuildImage(args *BuildImageArgs) *project.Job {
	imageResource := &project.JobResource{
		Name: args.Image,
	}

	preparedDir := &project.TaskOutput{
		Directory: "prepared",
	}

	prepareImageResource := &project.JobResource{
		Name:    args.Prepare.Name,
		Trigger: true,
	}

	taskPrepare := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "prepare",
		Image:    prepareImageResource,
		Run: &primitive.Location{
			Volume: &project.JobResource{
				Name:    ConcourseBuilderGitName,
				Trigger: true,
			},
			RelativePath: "scripts/docker_image_prepare.sh",
		},
		Environment: map[string]interface{}{
			"DOCKERFILE_DIR": args.DockerFileResource,
			"FROM_IMAGE":     (*image.FromParam)(args.From),
		},
		Outputs: []project.IOutput{
			preparedDir,
		},
	}

	if args.Eval != "" {
		taskPrepare.Environment["EVAL"] = args.Eval
	}

	imageSource := args.From.Source.(*image.Source)
	public := imageSource.Registry.Public()

	fromImageResource := &project.JobResource{
		Name:    args.From.Name,
		Trigger: true,
	}

	if !public {
		fromImageResource.GetParams = &resource.ImageGetParams{
			Save: true,
		}
	}

	putImage := &project.PutStep{
		JobResource: imageResource,
		Params: &image.PutParams{
			FromImage: fromImageResource,
			Load:      !public,
			Build: &primitive.Location{
				RelativePath: preparedDir.Path(),
			},
			BuildArgs: args.BuildArgs,
		},
		GetParams: &resource.ImageGetParams{
			SkipDownload: true,
		},
	}

	imageJob := &project.Job{
		Name: project.JobName(args.Name + "-image"),
		Groups: project.JobGroups{
			ImagesGroup,
		},
		Steps: args.PreprepareSteps,
	}

	imageJob.Steps = append(imageJob.Steps, taskPrepare, putImage)

	return imageJob
}
