package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type CurlImageJobArgs struct {
	LinuxImageResource  *project.Resource
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
}

func CurlImageJob(args *CurlImageJobArgs) *project.Resource {
	resourceName := project.ResourceName("curl-image")
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
			Repository: "concourse-builder/curl-image",
		},
	}

	if !needJob {
		return imageResource
	}

	steps := `RUN set -ex \
# install curl \
&& apt-get update \
&& apt-get install -y curl \
\
# cleanup \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*`

	job := BuildImage(
		&BuildImageArgs{
			ResourceRegistry:    args.ResourceRegistry,
			Prepare:             image.Ubuntu,
			From:                args.LinuxImageResource,
			Name:                "curl",
			DockerFileSteps:     steps,
			Image:               imageResource,
		})
	job.AddToGroup(project.SystemGroup)

	imageResource.NeedJobs(job)

	return imageResource
}
