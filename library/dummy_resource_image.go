package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type DummyResourceImageJobArgs struct {
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
}

func DummyResourceJob(args *DummyResourceImageJobArgs) *ResourceImageSource {
	resourceName := project.ResourceName("dummy_resource-image")
	imageResource := args.ResourceRegistry.GetResource(resourceName)
	if imageResource != nil {
		return (*ResourceImageSource)(imageResource)
	}

	imageResource = &project.Resource{
		Name:  resourceName,
		Type:  resource.ImageResourceType.Name,
		Scope: project.TeamScope,
		Source: &image.Source{
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/dummy_resource-image",
		},
	}

	args.ResourceRegistry.MustRegister(imageResource)

	dockerSteps := &primitive.Location{
		Volume:       args.ResourceRegistry.JobResource(args.ConcourseBuilderGit, true, nil),
		RelativePath: "docker/dummy_resource",
	}

	job := BuildImage(
		&BuildImageArgs{
			ResourceRegistry:   args.ResourceRegistry,
			PrepareImage:       image.Ubuntu,
			From:               image.Alpine,
			Name:               "dummy_resource",
			DockerFileResource: dockerSteps,
			Image:              imageResource,
		})
	job.AddToGroup(project.SystemGroup, project.ResourceTypeGroup)

	imageResource.NeedJobs(job)

	return (*ResourceImageSource)(imageResource)
}

func DummyResourceType(args *DummyResourceImageJobArgs) *project.ResourceType {
	source := DummyResourceJob(args)

	dummyResourceType := &project.ResourceType{
		Name:   "dummy",
		Type:   model.DockerImageType,
		Source: source,
	}

	project.GlobalTypeRegistry.MustRegisterType(dummyResourceType)

	return dummyResourceType
}
