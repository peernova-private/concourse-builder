package sdp

import (
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func GitImageJob() *project.Job {
	gitImageResource := &model.Resource{
		Name: "git-image",
		Type: resource.ImageResourceType.Name,
	}

	putGitImage := &project.PutStep{
		Resource: gitImageResource,
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
	return gitImageJob
}
