package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type CLangFormatImageJobArgs struct {
	ConcourseBuilderGitSource *GitSource
	ImageRegistry             *ImageRegistry
	ResourceRegistry          *project.ResourceRegistry
	Tag                       ImageTag
}

func CLangFormatImageJob(args *CLangFormatImageJobArgs) (*project.Resource, *project.Job) {
	resourceName := project.ResourceName("clang_format-image")
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
			Repository: "concourse-builder/clang_format-image",
		},
	}

	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/clang-format",
	}

	args.ResourceRegistry.MustRegister(UbuntuImage)

	job := BuildImage(
		UbuntuImage,
		UbuntuImage,
		&BuildImageArgs{
			Name:               "clang-format",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
		})
	job.AddToGroup(project.SystemGroup)

	image.NeededJobs = project.Jobs{job}
	args.ResourceRegistry.MustRegister(image)

	return image, job
}
