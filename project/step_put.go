package project

import (
	"time"

	"github.com/concourse-friends/concourse-builder/model"
)

type PutStep struct {
	// The resource that will be put
	Resource *model.Resource

	// How many attempts before give up
	Attempts int

	// Time duration for the get operation to timeout
	Timeout time.Duration

	// Additional resource specific parameters
	Params interface{}

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
