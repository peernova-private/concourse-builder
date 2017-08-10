package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

func FlyImageJob(concourse *library.Concourse, curlResource *project.Resource, flyImage project.ResourceName) *project.Job {
	dockerSteps := &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/fly_steps",
	}

	evalFlyVersion := "echo ENV FLY_VERSION=`curl " + concourse.URL + "/api/v1/info | " +
		"awk -F ',' ' { print $1 } ' | awk -F ':' ' { print $2 } '`"

	job := library.BuildImage(
		curlResource,
		curlResource,
		&library.BuildImageArgs{
			Name:               "fly",
			DockerFileResource: dockerSteps,
			Image:              flyImage,
			Eval:               evalFlyVersion,
		})

	job.AddToGroup(project.SystemGroup)
	return job
}
