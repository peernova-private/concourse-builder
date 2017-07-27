package library

import (
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

var imagesGroup = &project.JobGroup{
	Name: "images",
}

func BuildImage(name string, dockerFileResource project.IParamValue, image project.ResourceName) *project.Job {
	imageResource := &project.JobResource{
		Name: image,
	}

	ubuntuImageResource := &project.JobResource{
		Name:    UbuntuImage.Name,
		Trigger: true,
	}

	preparedDir := &TaskOutput{
		Directory: "prepared",
	}

	taskPrepare := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "prepare",
		Image:    ubuntuImageResource,
		Run: &Location{
			Volume: &project.JobResource{
				Name:    ConcourseBuilderGit,
				Trigger: true,
			},
			RelativePath: "scripts/docker_image_prepare.sh",
		},
		Params: map[string]interface{}{
			"DOCKER_STEPS": dockerFileResource,
			"FROM_IMAGE":   (*FromParam)(UbuntuImage),
		},
		Outputs: []project.IOutput{
			preparedDir,
		},
	}

	putImage := &project.PutStep{
		JobResource: imageResource,
		Params: &ImagePutParams{
			FromImage: ubuntuImageResource,
			Build: &Location{
				RelativePath: preparedDir.Path(),
			},
		},
		GetParams: &resource.ImageGetParams{
			SkipDownload: true,
		},
	}

	imageJob := &project.Job{
		Name: project.JobName(name + "-image"),
		Groups: project.JobGroups{
			imagesGroup,
		},
		Steps: project.ISteps{
			taskPrepare,
			putImage,
		},
	}
	return imageJob
}
