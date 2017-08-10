package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

func CurlImageJob(curlImage project.ResourceName) *project.Job {
	dockerSteps := &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/curl_steps",
	}

	job := library.BuildImage(
		library.UbuntuImage,
		library.UbuntuImage,
		&library.BuildImageArgs{
			Name:               "curl",
			DockerFileResource: dockerSteps,
			Image:              curlImage,
		})

	job.AddToGroup(project.SystemGroup)
	return job
}
