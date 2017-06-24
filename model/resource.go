package model

import (
	"time"
)

// A name of a resource
type ResourceName string

// Resource used from jobs
type Resource struct {
	// The resource name
	Name ResourceName `yaml:"type"`

	// The resource type by name
	Type ResourceTypeName `yaml:"type"`

	// Depending on the type a source that allows for the resource to be loaded
	Source interface{} `yaml:"type"`

	// Interval between checks for the resource
	CheckEvery time.Duration `yaml:"check_every",omitempty`
}

// A collection of resources
type Resources []*Resource
