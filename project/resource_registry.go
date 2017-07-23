package project

import (
	"fmt"

	"github.com/concourse-friends/concourse-builder/model"
)

type Resource struct {
	Name   ResourceName
	Type   model.ResourceTypeName
	Source IJobResourceSource
}

// An object that tracks collection of resources by name
type ResourceRegistry struct {
	resources map[ResourceName]*Resource
}

func NewResourceRegistry() *ResourceRegistry {
	return &ResourceRegistry{
		resources: make(map[ResourceName]*Resource),
	}
}

func (r *ResourceRegistry) MustRegister(resource *Resource) {
	_, ok := r.resources[resource.Name]
	if ok {
		// TODO: check if register more than once, the registered is the same
		return
	}

	r.resources[resource.Name] = resource
}

func (r *ResourceRegistry) MustGetResource(name ResourceName) *Resource {
	res, ok := r.resources[name]
	if !ok {
		panic(fmt.Sprintf("Resource %s is not found in the registry", name))
	}
	return res
}
