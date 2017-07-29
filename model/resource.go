package model

// A name of a resource
type ResourceName string

func (rn ResourceName) MarshalYAML() (interface{}, error) {
	return string(rn), nil
}

// Resource used from jobs
type Resource struct {
	// The resource name
	Name ResourceName

	// The resource type by name
	Type ResourceTypeName

	// Depending on the type a source that allows for the resource to be loaded
	Source interface{} `yaml:",omitempty"`

	// Interval between checks for the resource
	CheckEvery Duration `yaml:"check_every,omitempty"`
}

func (r *Resource) Path() string {
	return string(r.Name)
}

// A collection of resources
type Resources []*Resource
