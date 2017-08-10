package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
	"fmt"
)

func FlyImageJob(concourse *library.Concourse, curlResource *project.Resource, flyImage project.ResourceName) *project.Job {
	dockerSteps := &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/fly_steps",
	}
	var insecureArg string

	if concourse.Insecure {
		insecureArg = "-k"
	}
	evalFlyVersion := fmt.Sprintf("echo ENV FLY_VERSION=`curl %s/api/v1/info %s | " +
		"awk -F ',' ' { print $1 } ' | awk -F ':' ' { print $2 } '`", concourse.URL, insecureArg)

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
