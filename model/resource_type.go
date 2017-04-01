package model

// A name of a resource type
type ResourceTypeName string

// A type of a resource type (silly I know)
type ResourceTypeTypeName string

// Resource type needed to define the resources
type ResourceType struct {
	// The name of the resource type
	Name ResourceTypeName

	// The type of the resource type by name
	Type ResourceTypeTypeName

	// The source of the resource type
	Source interface{}
}

// A collection of resource types
type ResourceTypes []*ResourceType
