package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type GitImageJobArgs struct {
	*CurlImageJobArgs
}

func GitImageJob(args *GitImageJobArgs) (*project.Resource, *project.Job) {
	resourceName := project.ResourceName("git-image")
	image := args.ResourceRegistry.GetResource(resourceName)
	if image != nil {
		return image, image.NeededJobs[0]
	}

	curlImage, _ := CurlImageJob(args.CurlImageJobArgs)

	image = &project.Resource{
		Name: resourceName,
		Type: resource.ImageResourceType.Name,
		Source: &ImageSource{
			Tag:        args.Tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/git-image",
		},
	}

	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/git",
	}

	ubuntuImage := args.ResourceRegistry.MustRegister(UbuntuImage)

	job := BuildImage(
		ubuntuImage,
		curlImage,
		&BuildImageArgs{
			Name:               "git",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
		})

	image.NeededJobs = project.Jobs{job}
	image = args.ResourceRegistry.MustRegister(image)

	return image, job
}
