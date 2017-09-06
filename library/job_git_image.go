package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
	"github.com/jinzhu/copier"
)

type GitImageJobArgs struct {
	LinuxImageResource  *project.Resource
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
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

	tag, needJob := image.BranchImageTag(args.ConcourseBuilderGit.Source.(*GitSource).Branch)

	imageResource = &project.Resource{
		Name: resourceName,
		Type: resource.ImageResourceType.Name,
		Source: &image.Source{
			Tag:        tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/git-image",
		},
	}

	if !needJob {
		return imageResource
	}

	dockerSteps := &primitive.Location{
		Volume:       args.ResourceRegistry.JobResource(args.ConcourseBuilderGit, true, nil),
		RelativePath: "docker/git",
	}

	job := BuildImage(
		&BuildImageArgs{
			ConcourseBuilderGit: args.ConcourseBuilderGit,
			ResourceRegistry:    args.ResourceRegistry,
			Prepare:             image.Ubuntu,
			From:                curlImage,
			Name:                "git",
			DockerFileResource:  dockerSteps,
			Image:               imageResource,
		})

	imageResource.NeedJobs(job)

	return imageResource
}
