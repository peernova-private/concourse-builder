package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type CLangFormatImageJobArgs struct {
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
}

func CLangFormatImageJob(args *CLangFormatImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("clang_format-image")
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
			Repository: "concourse-builder/clang_format-image",
		},
	}

	if !needJob {
		return imageResource
	}

	dockerSteps := &primitive.Location{
		Volume:       args.ResourceRegistry.JobResource(args.ConcourseBuilderGit, true, nil),
		RelativePath: "docker/clang-format",
	}

	job := BuildImage(
		&BuildImageArgs{
			ConcourseBuilderGit: args.ConcourseBuilderGit,
			ResourceRegistry:   args.ResourceRegistry,
			Prepare:            image.Ubuntu,
			From:               image.Ubuntu,
			Name:               "clang-format",
			DockerFileResource: dockerSteps,
			Image:              imageResource,
		})
	job.AddToGroup(project.SystemGroup)

	imageResource.NeededJobs = project.Jobs{job}

	return imageResource
}
