package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

func GitImageJob(concourseBuilder, gitImage project.ResourceName) (*project.Job, error) {
	var concourseBuilderGirResource = &project.JobResource{
		Name:    concourseBuilder,
		Trigger: true,
	}

	gitImageResource := &project.JobResource{
		Name: gitImage,
	}

	putGitImage := &project.PutStep{
		JobResource: gitImageResource,
		Params: &library.ImagePutParams{
			Build: &library.Location{
				Volume: concourseBuilderGirResource,
				Path:   "template/sdp/docker/git-image",
			},
		},
	}

	gitImageJob := &project.Job{
		Name: "git-image",
		Groups: project.JobGroups{
			imagesGroup,
		},
		Steps: project.ISteps{
			putGitImage,
		},
	}
	return gitImageJob, nil
}
