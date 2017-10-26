package library

import (
	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/jinzhu/copier"
)

type SelfUpdateJobArgs struct {
	LinuxImageResource      *project.Resource
	ConcourseBuilderGit     *project.Resource
	ImageRegistry           *image.Registry
	GoImage                 *project.Resource
	PipelineResourceConfigImageJobArgs *project.Resource
	ResourceRegistry        *project.ResourceRegistry
	Concourse               *primitive.Concourse
	Environment             map[string]interface{}
	GenerateProjectLocation project.IRun
	Bucket                  *primitive.S3Bucket
}

func SelfUpdateJob(args *SelfUpdateJobArgs) (*project.Job, *project.Resource) {
	flyImageJobArgs := &FlyImageJobArgs{}
	copier.Copy(flyImageJobArgs, args)

	flyImage := FlyImageJob(flyImageJobArgs)

	flyImageResource := args.ResourceRegistry.JobResource(flyImage, true, nil)

	taskCheck := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "check",
		Image:    flyImageResource,
		Run: &primitive.Location{
			Volume: &primitive.Directory{
				Root: "/bin/fly",
			},
			RelativePath: "check_version.sh",
		},
		Environment: map[string]interface{}{},
	}

	args.Concourse.PublicAccessEnvironment(taskCheck.Environment)

	goImageResource := args.ResourceRegistry.JobResource(args.GoImage, true, nil)

	pipelinesDir := &project.TaskOutput{
		Directory: "pipelines",
	}

	taskPrepare := &project.TaskStep{
		Platform:    model.LinuxPlatform,
		Name:        "prepare pipelines",
		Image:       goImageResource,
		Run:         args.GenerateProjectLocation,
		Environment: args.Environment,
		Outputs: []project.IOutput{
			pipelinesDir,
		},
	}

	taskPrepare.Environment["PIPELINES"] = "pipelines"

	taskUpdate := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "update pipelines",
		Image:    flyImageResource,
		Run: &primitive.Location{
			Volume: &primitive.Directory{
				Root: "/bin/fly",
			},
			RelativePath: "set_pipelines.sh",
		},
		Environment: map[string]interface{}{
			"PIPELINES": &primitive.Location{
				Volume: pipelinesDir,
			},
		},
	}
	args.Concourse.Environment(taskUpdate.Environment)

	dummyResourceImageJobArgs := &DummyResourceImageJobArgs{}
	copier.Copy(dummyResourceImageJobArgs, args)

	dummyResourceType := DummyResourceType(dummyResourceImageJobArgs)

	pipelineResource := &project.Resource{
		Name: "pipeline",
		Type: dummyResourceType.Name,
	}

	args.ResourceRegistry.MustRegister(pipelineResource)

	pipelinePut := &project.PutStep{
		Resource: pipelineResource,
	}

	pipelineresourceconfigImageJobArgs := &PipelineResourceConfigImageJobArgs{}
	copier.Copy(pipelineresourceconfigImageJobArgs, args)

	pipelineresourceconfigImage := PipelineResourceConfigImageJob(pipelineresourceconfigImageJobArgs)

	pipelineresourceconfigImageResource := args.ResourceRegistry.JobResource((*project.Resource)(pipelineresourceconfigImage), true, nil)

	pipelineConfigCheck := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "check",
		Image:    pipelineresourceconfigImageResource,
		Run: &primitive.Location{
			Volume: &primitive.Directory{
				Root: "/opt/resource",
			},
			RelativePath: "check",
		},
		Environment: map[string]interface{}{},
	}

	args.Concourse.PublicAccessEnvironment(pipelineConfigCheck.Environment)
	/*pipelineResourceConfigType := PipelineResourceConfigType(pipelineresourceconfigImageJobArgs)

	pipelineConfig := &project.Resource{
		Name: "pipeline-config",
		Type:  pipelineResourceConfigType.Name,
	}

	args.ResourceRegistry.MustRegister(pipelineConfig)

	pipelineConfigPut :=&project.PutStep{
		Resource: pipelineConfig,
	}*/


	updateJob := &project.Job{
		Name:   project.JobName("self-update"),
		Groups: project.JobGroups{},
		Steps: project.ISteps{
			taskCheck,
			taskPrepare,
			taskUpdate,
			pipelinePut,
			pipelineConfigCheck,
		},
	}

	pipelineResource.NeedJobs(updateJob)

	return updateJob, pipelineResource
}
