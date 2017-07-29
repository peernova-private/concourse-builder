package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

func FlyImageJob(flyImage project.ResourceName) *project.Job {
	dockerSteps := &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGit,
			Trigger: true,
		},
		RelativePath: "docker/fly_steps",
	}

	return library.BuildImage(&library.BuildImageArgs{
		Name:               "fly",
		DockerFileResource: dockerSteps,
		Image:              flyImage,
		BuildArgs: map[string]interface{}{
			"FLY_VERSION": "v3.3.1",
		},
	})
}
