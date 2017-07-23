package project

import (
	"github.com/concourse-friends/concourse-builder/model"
)

type ResourceName string

type IJobResourceSource interface {
	ModelSource() interface{}
}

type JobResource struct {
	Name    ResourceName
	Trigger bool
}

type JobResources []*JobResource

func (jr *JobResource) Path() string {
	return string(jr.Name)
}

func (jr *JobResource) Model(registry *ResourceRegistry) *model.Resource {
	res := registry.MustGetResource(jr.Name)

	modelResource := &model.Resource{
		Name: model.ResourceName(jr.Name),
		Type: res.Type,
	}

	if res.Source != nil {
		modelResource.Source = res.Source.ModelSource()
	}

	return modelResource
}
