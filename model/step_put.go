package model

import (
	"time"
)

// A step for creating new version of a resource (or any related side effect)
type Put struct {
	// The resource that will be used by name
	Put ResourceName

	// How many attempts before give up
	Attempts int `yaml:",omitempty"`

	// Time duration for the get operation to timeout
	Timeout time.Duration `yaml:",omitempty"`

	// Additional resource specific parameters
	Params interface{} `yaml:",omitempty"`

	// Additional resource specific parameters for the get operation that will follow the put operation
	GetParams interface{} `yaml:"get_params,omitempty"`
}
