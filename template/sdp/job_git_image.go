package sdp

import (
	"path"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func GitImageJob(specification SdpSpecification) *project.Job {
	gitImageResource := &model.Resource{
		Name:   "git-image",
		Type:   resource.ImageResourceType.Name,
		Source: path.Join(specification.DeployImageRepository().Domain, "concourse-builder/git-image"),
	}

	putGitImage := &project.PutStep{
		Resource: gitImageResource,
		Params: &resource.ImagePutParams{
			Build: &model.Location{
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
	return gitImageJob
}
