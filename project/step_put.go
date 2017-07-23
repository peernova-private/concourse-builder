package project

import (
	"time"

	"github.com/concourse-friends/concourse-builder/model"
)

type IPutParams interface {
	InputResources() (model.Resources, error)
}

type PutStep struct {
	// The resource that will be put
	Resource *model.Resource

	// How many attempts before give up
	Attempts int

	// Time duration for the get operation to timeout
	Timeout time.Duration

	// Additional resource specific parameters
	Params IPutParams

	// Additional resource specific parameters for the get operation that will follow the put operation
	GetParams interface{}
}

func (ps *PutStep) Model() (model.IStep, error) {
	put := &model.Put{
		Put:       ps.Resource.Name,
		Attempts:  ps.Attempts,
		Timeout:   ps.Timeout,
		Params:    ps.Params,
		GetParams: ps.GetParams,
	}

	return put, nil
}

func (ps *PutStep) InputResources() (model.Resources, error) {
	if ps.Params != nil {
		return ps.Params.InputResources()
	}

	return nil, nil
}

func (ps *PutStep) OutputResource() (*model.Resource, error) {
	return ps.Resource, nil
}
