package model

// A name of a job
type JobName string

// A collection of job names
type JobNames []JobName

// Collection of steps to be executed together
type Job struct {
	// The name of the job
	Name JobName

	// Flag if the job is serial
	Serial bool `yaml:",omitempty"`

	// How many from the same job could run in parallel
	MaxInFlight int `yaml:"max_in_flight,omitempty"`

	// Serial groups the job belongs to
	SerialGroups SerialGroups `yaml:"serial_groups,omitempty"`

	// The execution plan
	Plan ISteps `yaml:",omitempty"`

	// A step to be executed on success
	OnSuccess IStep `yaml:"on_success,omitempty"`

	// A step to be executed on failure
	OnFailure IStep `yaml:"on_failure,omitempty"`
}

// Collection of jobs
type Jobs []*Job
