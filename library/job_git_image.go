package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func GitImageJob(imageRegistry *ImageRegistry, resourceRegistry *project.ResourceRegistry, tag ImageTag) (*project.Resource, *project.Job) {
	image := &project.Resource{
		Name: "git-image",
		Type: resource.ImageResourceType.Name,
		Source: &ImageSource{
			Tag:        tag,
			Registry:   imageRegistry,
			Repository: "concourse-builder/git-image",
		},
	}
	resourceRegistry.MustRegister(image)

	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/git",
	}

	job := BuildImage(
		UbuntuImage,
		UbuntuImage,
		&BuildImageArgs{
			Name:               "git",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
		})

	return image, job
}
