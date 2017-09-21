package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type GradleImageJobArgs struct {
	GradleImageResource *project.Resource
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
}

func GradleImageJob(args *GradleImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("root-gradle-image")
	imageResource := args.ResourceRegistry.GetResource(resourceName)
	if imageResource != nil {
		return imageResource
	}

	tag, needJob := image.BranchImageTag(args.ConcourseBuilderGit.Source.(*GitSource).Branch)

	imageResource = &project.Resource{
		Name: resourceName,
		Type: resource.ImageResourceType.Name,
		Source: &image.Source{
			Tag:        tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/gradle-image",
		},
	}

	if !needJob {
		return imageResource
	}

	steps := `USER root`

	job := BuildImage(
		&BuildImageArgs{
			ConcourseBuilderGit: args.ConcourseBuilderGit,
			ResourceRegistry:    args.ResourceRegistry,
			Prepare:             image.Ubuntu,
			From:                args.GradleImageResource,
			Name:                "gradle",
			DockerFileSteps:     steps,
			Image:               imageResource,
		})
	job.AddToGroup(project.SystemGroup)

	imageResource.NeedJobs(job)

	return imageResource
}
