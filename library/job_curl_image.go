package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type CurlImageJobArgs struct {
	ConcourseBuilderGitSource *GitSource
	ImageRegistry             *ImageRegistry
	ResourceRegistry          *project.ResourceRegistry
	Tag                       ImageTag
}

func CurlImageJob(args *CurlImageJobArgs) (*project.Resource, *project.Job) {
	resourceName := project.ResourceName("curl-image")
	image := args.ResourceRegistry.GetResource(resourceName)
	if image != nil {
		return image, image.NeededJobs[0]
	}

	image = &project.Resource{
		Name: resourceName,
		Type: resource.ImageResourceType.Name,
		Source: &ImageSource{
			Tag:        args.Tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/curl-image",
		},
	}

	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/curl",
	}

	ubuntuImage := args.ResourceRegistry.MustRegister(UbuntuImage)

	job := BuildImage(
		ubuntuImage,
		ubuntuImage,
		&BuildImageArgs{
			Name:               "curl",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
		})
	job.AddToGroup(project.SystemGroup)

	image.NeededJobs = project.Jobs{job}
	image = args.ResourceRegistry.MustRegister(image)

	return image, job
}
