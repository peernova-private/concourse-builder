package resource

import (
	"bytes"
	"fmt"

	"github.com/concourse-friends/concourse-builder/model"
	"gopkg.in/yaml.v2"
)

// An object that tracks collection of resources by name
type Registry struct {
	types     map[model.ResourceTypeName]*model.ResourceType
	resources map[model.ResourceName]*model.Resource
}

func (r *Registry) MustRegisterType(resourceType *model.ResourceType) {
	res, ok := r.types[resourceType.Name]
	if ok {
		current, err := yaml.Marshal(res)
		if err != nil {
			panic(err.Error())
		}
		new, err := yaml.Marshal(resourceType)
		if err != nil {
			panic(err.Error())
		}

		if bytes.Compare(current, new) != 0 {
			panic(fmt.Sprintf(
				"Resource type with name %s is already registered with different content", resourceType.Name))
		}
	}

	r.types[resourceType.Name] = resourceType
}

func (r *Registry) MustRegisterResource(resource *model.Resource) {
	res, ok := r.resources[resource.Name]
	if ok {
		current, err := yaml.Marshal(res)
		if err != nil {
			panic(err.Error())
		}
		new, err := yaml.Marshal(resource)
		if err != nil {
			panic(err.Error())
		}

		if bytes.Compare(current, new) != 0 {
			panic(fmt.Sprintf(
				"Resource with name %s is already registered with different content", resource.Name))
		}
	}

	r.resources[resource.Name] = resource
}

var GlobalRegistry = initRegistry()

func initRegistry() *Registry {
	registry := &Registry{
		types:     make(map[model.ResourceTypeName]*model.ResourceType),
		resources: make(map[model.ResourceName]*model.Resource),
	}

	return registry
}
