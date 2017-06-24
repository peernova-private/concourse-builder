package model

// A pipeline structure
type Pipeline struct {
	// Organized jobs
	Groups Groups `yaml:",omitempty"`

	// Resource types needed for defining the resources
	ResourceTypes ResourceTypes `yaml:"resource_types,omitempty"`

	// Resources needed for defining the jobs
	Resources Resources `yaml:",omitempty"`

	// Jobs in the pipeline
	Jobs Jobs `yaml:",omitempty"`
}
