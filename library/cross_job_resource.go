package library

import (
	"fmt"

	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func CrossResource(
	bucket *primitive.S3Bucket,
	job *project.Job,
	output *project.TaskOutput) *project.Resource {

	taskIndex, outputIndex := job.TaskOutputIndex(output)

	name := fmt.Sprintf("cross_%d_%d", taskIndex, outputIndex)

	resource := &project.Resource{
		Name: project.ResourceName(name + "-s3"),
		Type: resource.S3ResourceType.Name,
		Source: &S3Source{
			Bucker:        bucket,
			VersionedFile: name,
		},
	}

	putS3 := &project.PutStep{
		Resource: resource,
	}

	job.Steps = append(job.Steps, putS3)

	resource.NeedJobs(job)

	return resource
}
