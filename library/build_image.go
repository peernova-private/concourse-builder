package library

import (
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
	Name               string
	DockerFileResource project.IParamValue
	Image              project.ResourceName
	BuildArgs          map[string]interface{}
	Eval               string
}

func BuildImage(prepare *project.Resource, from *project.Resource, args *BuildImageArgs) *project.Job {
	imageResource := &project.JobResource{
		Name: args.Image,
	}

	preparedDir := &TaskOutput{
		Directory: "prepared",
	}

	prepareImageResource := &project.JobResource{
		Name:    prepare.Name,
		Trigger: true,
	}

	taskPrepare := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "prepare",
		Image:    prepareImageResource,
		Run: &Location{
			Volume: &project.JobResource{
				Name:    ConcourseBuilderGitName,
				Trigger: true,
			},
			RelativePath: "scripts/docker_image_prepare.sh",
		},
		Params: map[string]interface{}{
			"DOCKERFILE_DIR": args.DockerFileResource,
			"FROM_IMAGE":     (*FromParam)(from),
		},
		Outputs: []project.IOutput{
			preparedDir,
		},
	}

	if args.Eval != "" {
		taskPrepare.Params["EVAL"] = args.Eval
	}

	imageSource := from.Source.(*ImageSource)
	public := imageSource.Registry.Public()

	fromImageResource := &project.JobResource{
		Name:    from.Name,
		Trigger: true,
	}

	if !public {
		fromImageResource.GetParams = &resource.ImageGetParams{
			Save: true,
		}
	}

	putImage := &project.PutStep{
		JobResource: imageResource,
		Params: &ImagePutParams{
			FromImage: fromImageResource,
			Load:      !public,
			Build: &Location{
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
		Steps: project.ISteps{
			taskPrepare,
			putImage,
		},
	}
	return imageJob
}
