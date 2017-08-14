package library

import (
	"fmt"

	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type FlyImageJobArgs struct {
	*CurlImageJobArgs
	Concourse *Concourse
}

func FlyImageJob(args *FlyImageJobArgs) (*project.Resource, *project.Job) {
	resourceName := project.ResourceName("fly-image")
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
			Repository: "concourse-builder/fly-image",
		},
	}

	RegisterConcourseBuilderGit(args.ResourceRegistry, args.ConcourseBuilderGitSource)

	dockerSteps := &Location{
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
	args.ResourceRegistry.MustRegister(image)

	return image, job
}
