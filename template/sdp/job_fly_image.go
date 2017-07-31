package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

func FlyImageJob(flyVersion string, flyImage project.ResourceName) *project.Job {
	dockerSteps := &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/fly_steps",
	}

	job := library.BuildImage(&library.BuildImageArgs{
		Name:               "fly",
		DockerFileResource: dockerSteps,
		Image:              flyImage,
		BuildArgs: map[string]interface{}{
			"FLY_VERSION": flyVersion,
		},
	})

	job.AddToGroup(project.SystemGroup)
	return job
}
