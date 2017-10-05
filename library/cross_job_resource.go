package library

import (
	"fmt"
	"path"

	"github.com/concourse-friends/concourse-builder/library/image"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

func CrossResource(
	register *project.ResourceRegistry,
	bucket *primitive.S3Bucket,
	job *project.Job,
	output *project.TaskOutput) *project.JobResource {

	taskIndex, outputIndex := job.TaskOutputIndex(output)

	name := fmt.Sprintf("%s_cross_%d_%d", job.Name, taskIndex, outputIndex)

	s3 := &project.Resource{
		Name:  project.ConvertToResourceName(name + "-s3"),
		Type:  resource.S3ResourceType.Name,
		Scope: project.UniversalScope,
		Source: &S3Source{
			Bucker:        bucket,
			VersionedFile: name + ".tar.gz",
		},
	}

	archive := &project.TaskOutput{
		Directory: "archive",
	}

	taskZip := &project.TaskStep{
		Platform: model.LinuxPlatform,
		Name:     "compress",
		Image:    register.JobResource(image.Ubuntu, true, nil),
		Run: &primitive.Location{
			Volume: &primitive.Directory{
				Root: "/bin",
			},
			RelativePath: "tar",
		},
		Arguments: []interface{}{
			"-zcvf", path.Join(archive.Path(), "archive.tar.gz"),
			"-C", &primitive.Location{
				Volume: output,
			},
			".",
		},
		Outputs: []project.IOutput{
			archive,
		},
	}

	putS3 := &project.PutStep{
		Resource: s3,
		Params: &S3PutParams{
			File: path.Join(archive.Path(), "archive.tar.gz"),
		},
	}

	job.Steps = append(job.Steps, taskZip, putS3)

	s3.NeedJobs(job)

	jobResurce := register.JobResource(s3, true, &resource.S3GetParams{
		Unpack: true,
	})

	jobResurce.PreferredPath = output.Directory

	return jobResurce
}
