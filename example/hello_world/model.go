package hello_world

import "github.com/concourse-friends/concourse-builder/model"

var Pipeline = model.Pipeline{
	Jobs: model.Jobs{
		{
			Name: "hello-world",
			Plan: []model.IStep{
				model.Task{
					Task: "say-hello",
					Config: &model.TaskConfig{
						Platform: "linux",
						ImageResource: model.TaskImageResource{
							Type: model.TaskImageResourceTypeDocker,
							Source: model.TaskImageResourceDockerSource{
								Repository: "ubuntu",
							},
						},
						Run: model.TaskRun{
							Path: "echo",
							Args: []string{
								"Hello, world!",
							},
						},
					},
				},
			},
		},
	},
}
