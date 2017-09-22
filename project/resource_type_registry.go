package project

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v2"
)

// An object that tracks collection of resources by name
type TypeRegistry struct {
	types map[ResourceTypeName]*ResourceType
}

func (r *TypeRegistry) MustRegisterType(resourceType *ResourceType) {
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
		return
	}

	r.types[resourceType.Name] = resourceType
}

func (r *TypeRegistry) RegisterType(resourceTypeName ResourceTypeName) *ResourceType {
	return r.types[resourceTypeName]
}

var GlobalTypeRegistry = initTypeRegistry()

func initTypeRegistry() *TypeRegistry {
	registry := &TypeRegistry{
		types: make(map[ResourceTypeName]*ResourceType),
	}

	return registry
}
