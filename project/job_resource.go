package project

import "github.com/concourse-friends/concourse-builder/model"

type IJobResourceSource interface {
	ModelSource() interface{}
}

type JobResource struct {
	Name    model.ResourceName
	Type    model.ResourceTypeName
	Source  IJobResourceSource
	Trigger bool
}

type JobResources []*JobResource

func (jr *JobResource) Path() string {
	return string(jr.Name)
}

func (jr *JobResource) Model() *model.Resource {
	modelResource := &model.Resource{
		Name: jr.Name,
		Type: jr.Type,
	}

	if jr.Source != nil {
		modelResource.Source = jr.Source.ModelSource()
	}

	return modelResource
}
