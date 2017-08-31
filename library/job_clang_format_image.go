package library

import (
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type CLangFormatImageJobArgs struct {
	ConcourseBuilderGitSource *GitSource
	ImageRegistry             *ImageRegistry
	ResourceRegistry          *project.ResourceRegistry
}

func CLangFormatImageJob(args *CLangFormatImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("clang_format-image")
	image := args.ResourceRegistry.GetResource(resourceName)
	if image != nil {
		return image
	}

	tag, needJob := BranchImageTag(args.ConcourseBuilderGitSource.Branch)

	image = &project.Resource{
		Name: resourceName,
		Type: resource.ImageResourceType.Name,
		Source: &ImageSource{
			Tag:        tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/clang_format-image",
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

	return image
}
