package model

import (
	"time"
)

// A step for obtaining a resource
type Get struct {
	// Resource to get by name
	Get ResourceName `yaml:",omitempty"`

	// How many attempts before give up
	Attempts int `yaml:",omitempty"`

	// Is change in the resource is a trigger for the job this get belongs too
	Trigger bool `yaml:",omitempty"`

	// Time duration for the get operation to timeout
	Timeout time.Duration `yaml:",omitempty"`

	// Version selection strategy
	Version string `yaml:",omitempty"`

	// Which jobs validate this resource by names
	Passed JobNames `yaml:",omitempty"`

	// Additional resource specific parameters
	Params interface{} `yaml:",omitempty"`
}
