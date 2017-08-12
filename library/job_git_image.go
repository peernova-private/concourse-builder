package library

import (
	"github.com/concourse-friends/concourse-builder/project"
)

func GitImageJob(gitImage project.ResourceName) *project.Job {
	dockerSteps := &Location{
		Volume: &project.JobResource{
			Name:    ConcourseBuilderGitName,
			Trigger: true,
		},
		RelativePath: "docker/git",
	}

	return BuildImage(
		UbuntuImage,
		UbuntuImage,
		&BuildImageArgs{
			Name:               "git",
			DockerFileResource: dockerSteps,
			Image:              gitImage,
		})
}
