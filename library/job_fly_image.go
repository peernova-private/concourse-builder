package library

import (
	"fmt"

	"github.com/concourse-friends/concourse-builder/project"
)

func FlyImageJob(concourse *Concourse, curlResource *project.Resource, flyImage project.ResourceName) *project.Job {
	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/fly",
	}
	var insecureArg string

	if concourse.Insecure {
		insecureArg = " -k"
	}
	evalFlyVersion := fmt.Sprintf("echo ENV FLY_VERSION=`curl %s/api/v1/info%s | "+
		"awk -F ',' ' { print $1 } ' | awk -F ':' ' { print $2 } '`", concourse.URL, insecureArg)

	job := BuildImage(
		curlResource,
		curlResource,
		&BuildImageArgs{
			Name:               "fly",
			DockerFileResource: dockerSteps,
			Image:              flyImage,
			Eval:               evalFlyVersion,
		})

	job.AddToGroup(project.SystemGroup)
	return job
}
