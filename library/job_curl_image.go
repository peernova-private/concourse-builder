package library

import (
	"github.com/concourse-friends/concourse-builder/project"
)

func CurlImageJob(curlImage project.ResourceName) *project.Job {
	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/curl",
	}

	job := BuildImage(
		UbuntuImage,
		UbuntuImage,
		&BuildImageArgs{
			Name:               "curl",
			DockerFileResource: dockerSteps,
			Image:              curlImage,
		})

	job.AddToGroup(project.SystemGroup)
	return job
}
