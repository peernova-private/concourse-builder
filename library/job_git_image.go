package library

import (
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
	"github.com/jinzhu/copier"
)

type GitImageJobArgs struct {
	ConcourseBuilderGitSource *GitSource
	ImageRegistry             *ImageRegistry
	ResourceRegistry          *project.ResourceRegistry
}

func GitImageJob(args *GitImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("git-image")
	image := args.ResourceRegistry.GetResource(resourceName)
	if image != nil {
		return image
	}

	curlImageJobArgs := &CurlImageJobArgs{}
	copier.Copy(curlImageJobArgs, args)

	curlImage := CurlImageJob(curlImageJobArgs)

	tag, needJob := BranchImageTag(args.ConcourseBuilderGitSource.Branch)

	image = &project.Resource{
		Name: resourceName,
		Type: resource.ImageResourceType.Name,
		Source: &ImageSource{
			Tag:        tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/git-image",
		},
	}
	args.ResourceRegistry.MustRegister(image)

	if !needJob {
		return image
	}

	RegisterConcourseBuilderGit(args.ResourceRegistry, args.ConcourseBuilderGitSource)

	dockerSteps := &primitive.Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/git",
	}

	args.ResourceRegistry.MustRegister(UbuntuImage)

	job := BuildImage(
		UbuntuImage,
		curlImage,
		&BuildImageArgs{
			Name:               "git",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
		})

	image.NeededJobs = project.Jobs{job}

	return image
}
