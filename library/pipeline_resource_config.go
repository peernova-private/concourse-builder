// Copyright (C) 2017 PeerNova, Inc.
//
// All rights reserved.
//
// PeerNova and Cuneiform are trademarks of PeerNova, Inc. References to
// third-party marks or brands are the property of their respective owners.
// No rights or licenses are granted, express or implied, unless set forth in
// a written agreement signed by PeerNova, Inc. You may not distribute,
// disseminate, copy, record, modify, enhance, supplement, create derivative
// works from, adapt, or translate any content contained herein except as
// otherwise expressly permitted pursuant to a written agreement signed by
// PeerNova, Inc.

package library

import (
"github.com/concourse-friends/concourse-builder/library/image"
"github.com/concourse-friends/concourse-builder/library/primitive"
"github.com/concourse-friends/concourse-builder/project"
"github.com/concourse-friends/concourse-builder/resource"
	"github.com/concourse-friends/concourse-builder/model"
)

type PipelineResourceConfigImageJobArgs struct {
	LinuxImageResource  *project.Resource
	ConcourseBuilderGit *project.Resource
	ImageRegistry       *image.Registry
	ResourceRegistry    *project.ResourceRegistry
	Concourse           *primitive.Concourse
}

func PipelineResourceConfigImageJob(args *PipelineResourceConfigImageJobArgs) *ResourceImageSource {
	resourceName := project.ResourceName("pipeline-resource-config")
	imageResource := args.ResourceRegistry.GetResource(resourceName)

	if imageResource != nil {
		return (*ResourceImageSource)(imageResource)
	}

	imageResource = &project.Resource{
		Name:  resourceName,
		Type:  resource.ImageResourceType.Name,
		Scope: project.AllTeamsScope,
		Source: &image.Source{
			Registry:   args.ImageRegistry,
			Repository: "concourse-builder/pipeline-resource",
		},
	}

	args.ResourceRegistry.MustRegister(imageResource)
	dockerSteps := &primitive.Location{
		Volume:       args.ResourceRegistry.JobResource(args.ConcourseBuilderGit, true, nil),
		RelativePath: "docker/pipeline_resource",
	}

	job := BuildImage(
		&BuildImageArgs{
			ResourceRegistry:   args.ResourceRegistry,
			PrepareImage:       image.Ubuntu,
			From:               args.LinuxImageResource,
			Name:               "pipeline-resource",
			DockerFileResource: dockerSteps,
			Image:              imageResource,
		})
	job.AddToGroup(project.SystemGroup)

	imageResource.NeedJobs(job)

	return (*ResourceImageSource)(imageResource)
}


func PipelineResourceConfigType(args *PipelineResourceConfigImageJobArgs) *project.ResourceType {
	source := PipelineResourceConfigImageJob(args)

	pipelineConfigResourceType := &project.ResourceType{
		Name:   "pipeline-resource",
		Type:   model.DockerImageType,
		Source: source,
	}

	project.GlobalTypeRegistry.MustRegisterType(pipelineConfigResourceType)

	return pipelineConfigResourceType
}
