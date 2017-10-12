package project

import (
	"github.com/concourse-friends/concourse-builder/model"
)

var ResourceTypeGroup = &JobGroup{
	Name: "res_types",
	Before: JobGroups{
		SystemGroup,
	},
}

type ResourceTypeName string

type IResourceTypeSource interface {
	ResourceName() ResourceName
	ResourceScope() Scope
	ModelSource(scope Scope, info *ScopeInfo) interface{}
	NeededJobs() Jobs
}

type ResourceType struct {
	Name   ResourceTypeName
	Type   model.ResourceTypeTypeName
	Source IResourceTypeSource
}

type ResourceTypes []*ResourceType

func (rt *ResourceType) IsSystem() bool {
	return string(rt.Type) == string(model.SystemResourceTypeName)
}

func (rt *ResourceType) Model(info *ScopeInfo) *model.ResourceType {
	resourceType := &model.ResourceType{
		Name: model.ResourceTypeName(rt.Name),
		Type: rt.Type,
	}

	if rt.Source != nil {
		resourceType.Source = rt.Source.ModelSource(rt.Source.ResourceScope(), info)
	}

	return resourceType
}
