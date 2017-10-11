package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type GradleImageJobArgs struct {
	GradleImageResource *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
}

func GradleImageJob(args *GradleImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("root-gradle-image")
	imageResource := args.ResourceRegistry.GetResource(resourceName)
	if imageResource != nil {
		return imageResource
	}

	imageResource = &project.Resource{
		Name:  resourceName,
		Type:  resource.ImageResourceType.Name,
		Scope: project.TeamScope,
		Source: &image.Source{
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/gradle-image",
		},
	}

	steps := `USER root`

	job := BuildImage(
		&BuildImageArgs{
			ResourceRegistry: args.ResourceRegistry,
			PrepareImage:     image.Ubuntu,
			From:             args.GradleImageResource,
			Name:             "gradle",
			DockerFileSteps:  steps,
			Image:            imageResource,
		})
	job.AddToGroup(project.SystemGroup)

	imageResource.NeedJobs(job)

	return imageResource
}
