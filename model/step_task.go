package model

import (
	"time"
)

// Type of the image the task will use to run on
type TaskImageResourceType string

// A docker image type, commonly used
var TaskImageResourceTypeDocker = TaskImageResourceType("docker-image")

// A source of a image resource
type TaskImageResourceDockerSource struct {
	// The repository of the docker image
	Repository string

	// The tag of the image
	Tag string `yaml:",omitempty"`
}

// A task image resource
type TaskImageResource struct {
	// The type of the task image
	Type TaskImageResourceType

	// Type specific source
	Source interface{}
}

// An input for a task
type TaskInput struct {
	// The name of the resource to input
	Name ResourceName `yaml:",omitempty"`

	// Location where the resource to be obtained
	Path string `yaml:",omitempty"`
}

// A run for a task
type TaskRun struct {
	// The path to the executable
	Path string `yaml:",omitempty"`

	// Arguments to the executable
	Args []string `yaml:",omitempty"`

	// Relative location for the execution
	Dir string `yaml:",omitempty"`

	// An user to use to execute the command
	User string `yaml:",omitempty"`
}

// An output for a task
type TaskOutput struct {
	// A name of a transfer resource to output
	Name string `yaml:",omitempty"`

	// The location of the artifacts
	Path string `yaml:",omitempty"`
}

// Configuration of a task
type TaskConfig struct {
	// What platform the task runs on
	Platform Platform `yaml:",omitempty"`

	// Image resource the task to use for container
	ImageResource TaskImageResource `yaml:"image_resource,omitempty"`

	// List of inputs
	Inputs []*TaskInput `yaml:",omitempty"`

	// Additional params (set as environment variables)
	Params map[string]interface{} `yaml:",omitempty"`

	// What to run
	Run *TaskRun `yaml:",omitempty"`

	// List of outputs from the task
	Outputs []*TaskOutput `yaml:",omitempty"`
}

// A name of a task
type TaskName string

// A Task to execute in a job
type Task struct {
	// The name of the task
	Task TaskName

	// A time duration in which the task to timeout
	Timeout time.Duration `yaml:",omitempty"`

	// Alternative to a Image resource from the config, not recommended from the concourse team,
	// but in fact in many cases preferable. The config image resource is not versioned, running the
	// same task over time may produce different results based on the used image resource. In order
	// Image to works it needs to be 'Get' first which makes it part of the visioning system od concourse.
	Image ResourceName `yaml:",omitempty"`

	// Optional. Default false. If set to true, the task will run with full capabilities,
	// as determined by the Garden backend the task runs on. For Linux-based backends it
	// typically determines whether or not the container will run in a separate user namespace,
	// and whether the root user is "actual" root (if set to true) or a user namespaced root
	// (if set to false, the default).
	Privileged bool `yaml:",omitempty"`

	// The configuration of the task
	Config *TaskConfig `yaml:",omitempty"`

	// Sub step that will be executed at the end of the task if it fail or not
	Ensure IStep `yaml:",omitempty"`

	// Sub step that will be executed if the task succeeds
	OnSuccess IStep `yaml:"on_success,omitempty"`

	// Sub step that will be executed if the task fail
	OnFailure IStep `yaml:"on_failure,omitempty"`

	// A number of attempts before the task is considered to fail
	Attempts int `yaml:",omitempty"`
}
