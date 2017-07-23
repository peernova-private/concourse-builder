package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func GitImageJob(specification SdpSpecification) (*project.Job, error) {
	privateKey, err := specification.GitPrivateKey()
	if err != nil {
		return nil, err
	}

	var concourseBuilderGirResource = &project.JobResource{
		Name: "concourse-builder-git",
		Type: resource.GitResourceType.Name,
		Source: &library.GitSource{
			URI:        "git@github.com:concourse-friends/concourse-builder.git",
			Branch:     "master",
			PrivateKey: privateKey,
		},
		Trigger: true,
	}

	gitImageResource := &project.JobResource{
		Name: "git-image",
		Type: resource.ImageResourceType.Name,
		Source: &library.ImageSource{
			Repository: specification.DeployImageRepository(),
			Location:   "concourse-builder/git-image",
		},
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
