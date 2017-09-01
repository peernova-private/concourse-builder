package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
	"github.com/jinzhu/copier"
)

type GitImageJobArgs struct {
	ConcourseBuilderGitSource *GitSource
	ImageRegistry             *image.Registry
	ResourceRegistry          *project.ResourceRegistry
}

func GitImageJob(args *GitImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("git-image")
	imageResource := args.ResourceRegistry.GetResource(resourceName)
	if imageResource != nil {
		return imageResource
	}

	curlImageJobArgs := &CurlImageJobArgs{}
	copier.Copy(curlImageJobArgs, args)

	curlImage := CurlImageJob(curlImageJobArgs)

	tag, needJob := image.BranchImageTag(args.ConcourseBuilderGitSource.Branch)

	imageResource = &project.Resource{
		Name: resourceName,
		Type: resource.ImageResourceType.Name,
		Source: &image.Source{
			Tag:        tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/git-image",
		},
	}
	args.ResourceRegistry.MustRegister(imageResource)

	if !needJob {
		return imageResource
	}

	RegisterConcourseBuilderGit(args.ResourceRegistry, args.ConcourseBuilderGitSource)

	dockerSteps := &primitive.Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/git",
	}

	args.ResourceRegistry.MustRegister(image.Ubuntu)

	job := BuildImage(
		image.Ubuntu,
		curlImage,
		&BuildImageArgs{
			Name:               "git",
			DockerFileResource: dockerSteps,
			Image:              imageResource.Name,
		})

	imageResource.NeededJobs = project.Jobs{job}

	return imageResource
}
