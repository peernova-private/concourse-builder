package library

import (
	"fmt"

	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type FlyImageJobArgs struct {
	ImageRegistry    *ImageRegistry
	ResourceRegistry *project.ResourceRegistry
	Tag              ImageTag
	Concourse        *Concourse
	CurlResource     *project.Resource
}

func FlyImageJob(args *FlyImageJobArgs) (*project.Resource, *project.Job) {
	image := &project.Resource{
		Name: "fly-image",
		Type: resource.ImageResourceType.Name,
		Source: &ImageSource{
			Tag:        args.Tag,
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/fly-image",
		},
	}
	args.ResourceRegistry.MustRegister(image)

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
		args.CurlResource,
		args.CurlResource,
		&BuildImageArgs{
			Name:               "fly",
			DockerFileResource: dockerSteps,
			Image:              image.Name,
			Eval:               evalFlyVersion,
		})

	job.AddToGroup(project.SystemGroup)
	return image, job
}
