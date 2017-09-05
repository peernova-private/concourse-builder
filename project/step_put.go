package project

import (
	"github.com/concourse-friends/concourse-builder/model"
)

type IPutParams interface {
	ModelParams() interface{}
}

type PutStep struct {
	// The resource that will be put
	Resource *Resource

	// Additional resource specific parameters
	Params IPutParams

	// Additional resource specific parameters for the get operation that will follow the put operation
	GetParams interface{}
}

func (ps *PutStep) Model() (model.IStep, error) {
	put := &model.Put{
		Put:       model.ResourceName(ps.Resource.Name),
		GetParams: ps.GetParams,
	}

	if ps.Params != nil {
		put.Params = ps.Params.ModelParams()
	}

	return put, nil
}

func (ps *PutStep) InputResources() (JobResources, error) {
	var resources JobResources

	if ps.Params != nil {
		if res, ok := ps.Params.(IInputResource); ok {
			resources = append(resources, res.InputResources()...)
		}
	}

	return resources, nil
}

func (ps *PutStep) OutputResource() (*Resource, error) {
	return ps.Resource, nil
}
