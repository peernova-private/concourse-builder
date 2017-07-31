package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

func GitImageJob(gitImage project.ResourceName) *project.Job {
	dockerSteps := &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/git_steps",
	}

	return library.BuildImage(&library.BuildImageArgs{
		Name:               "git",
		DockerFileResource: dockerSteps,
		Image:              gitImage,
	})
}
