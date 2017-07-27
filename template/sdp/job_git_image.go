package sdp

import (
	"github.com/concourse-friends/concourse-builder/library"
	"github.com/concourse-friends/concourse-builder/project"
)

func GitImageJob(concourseBuilder, gitImage project.ResourceName) *project.Job {
	dockerSteps := &library.Location{
		Volume: &project.JobResource{
			Name:    library.ConcourseBuilderGit,
			Trigger: true,
		},
		RelativePath: "docker/git_steps",
	}

	return library.BuildImage("git", dockerSteps, gitImage)
}
