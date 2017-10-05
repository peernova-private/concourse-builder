package library

import (
	"fmt"

	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
	"github.com/jinzhu/copier"
)

type FlyImageJobArgs struct {
	LinuxImageResource  *project.Resource
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
	Concourse           *primitive.Concourse
}

func FlyImageJob(args *FlyImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("fly-image")
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
			Repository: "concourse-builder/fly-image",
		},
	}

	if !needJob {
		return imageResource
	}

	dockerSteps := &primitive.Location{
		Volume:       args.ResourceRegistry.JobResource(args.ConcourseBuilderGit, true, nil),
		RelativePath: "docker/fly",
	}
	var insecureArg string

	if args.Concourse.Insecure {
		insecureArg = " -k"
	}
	evalFlyVersion := fmt.Sprintf("echo ENV FLY_VERSION=`curl %s/api/v1/info%s | "+
		"awk -F ',' ' { print $1 } ' | awk -F ':' ' { print $2 } '`", args.Concourse.URL, insecureArg)

	job := BuildImage(
		&BuildImageArgs{
			ResourceRegistry:   args.ResourceRegistry,
			PrepareImage:       curlImage,
			From:               curlImage,
			Name:               "fly",
			DockerFileResource: dockerSteps,
			Image:              imageResource,
			Eval:               evalFlyVersion,
		})
	job.AddToGroup(project.SystemGroup)

	imageResource.NeedJobs(job)

	return imageResource
}
