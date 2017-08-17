package project

// An object that tracks collection of resources by name
type ResourceRegistry struct {
	cross     map[ResourceName]ResourceHash
	resources map[ResourceHash]*Resource
}

func NewResourceRegistry() *ResourceRegistry {
	return &ResourceRegistry{
		cross:     make(map[ResourceName]ResourceHash),
		resources: make(map[ResourceHash]*Resource),
	}
}

func (rr *ResourceRegistry) MustRegister(resource *Resource) {
	if hash, ok := rr.cross[resource.Name]; ok {
		resource.Name = rr.resources[hash].Name
		return
	}

	hash := resource.MustHash()
	rr.cross[resource.Name] = hash

	if res, ok := rr.resources[hash]; ok {
		resource.Name = res.Name
		return
	}

	rr.resources[hash] = resource
}

func (rr *ResourceRegistry) GetResource(name ResourceName) *Resource {
	if hash, ok := rr.cross[name]; ok {
		if res, ok := rr.resources[hash]; ok {
			return res
		}
	}
	return nil
}

func (rr *ResourceRegistry) MustGetResource(name ResourceName) *Resource {
	return rr.resources[rr.cross[name]]
}
