package library

import (
	"fmt"

	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
	"github.com/jinzhu/copier"
)

type FlyImageJobArgs struct {
	ConcourseBuilderGitSource *GitSource
	ImageRegistry             *ImageRegistry
	ResourceRegistry          *project.ResourceRegistry
	Concourse                 *primitive.Concourse
}

func FlyImageJob(args *FlyImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("fly-image")
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
			Repository: "concourse-builder/fly-image",
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
		RelativePath: "docker/fly",
	}
	var insecureArg string

	if args.Concourse.Insecure {
		insecureArg = " -k"
	}
	evalFlyVersion := fmt.Sprintf("echo ENV FLY_VERSION=`curl %s/api/v1/info%s | "+
		"awk -F ',' ' { print $1 } ' | awk -F ':' ' { print $2 } '`", args.Concourse.URL, insecureArg)

	job := BuildImage(
		curlImage, // We need curlImage for prepare for eval to works
		curlImage,
		&BuildImageArgs{
			Name:               "fly",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
			Eval:               evalFlyVersion,
		})
	job.AddToGroup(project.SystemGroup)

	image.NeededJobs = project.Jobs{job}

	return image
}
