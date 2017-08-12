package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func CurlImageJob(imageRegistry *ImageRegistry, resourceRegistry *project.ResourceRegistry, tag ImageTag) (*project.Resource, *project.Job) {
	image := &project.Resource{
		Name: "curl-image",
		Type: resource.ImageResourceType.Name,
		Source: &ImageSource{
			Tag:        tag,
			Registry:   imageRegistry,
			Repository: "concourse-builder/curl-image",
		},
	}
	resourceRegistry.MustRegister(image)

	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/curl",
	}

	job := BuildImage(
		UbuntuImage,
		UbuntuImage,
		&BuildImageArgs{
			Name:               "curl",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
		})

	job.AddToGroup(project.SystemGroup)
	return image, job
}
