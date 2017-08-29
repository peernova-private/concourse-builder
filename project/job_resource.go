package project

import (
	"sort"
	"strings"

	"github.com/concourse-friends/concourse-builder/model"
)

type ResourceName string

func ConvertToResourceName(raw string) ResourceName {
	name := strings.Replace(raw, "/", "_", -1)
	return ResourceName(name)
}

type IJobResourceSource interface {
	ModelSource() interface{}
}

type JobResource struct {
	Name      ResourceName
	Trigger   bool
	GetParams interface{}
}

func (jr *JobResource) Path() string {
	return string(jr.Name)
}

func (jr *JobResource) Model(registry *ResourceRegistry) *model.Resource {
	res := registry.MustGetResource(jr.Name)

	modelResource := &model.Resource{
		Name:       model.ResourceName(jr.Name),
		Type:       res.Type,
		CheckEvery: res.CheckInterval,
	}

	if res.Source != nil {
		modelResource.Source = res.Source.ModelSource()
	}

	return modelResource
}

type JobResources []*JobResource

func (jr JobResources) Len() int {
	return len(jr)
}

func (jr JobResources) Swap(i, j int) {
	jr[i], jr[j] = jr[j], jr[i]
}

func (jr JobResources) Less(i, j int) bool {
	return jr[i].Name < jr[j].Name
}

func (jr JobResources) Deduplicate() JobResources {
	if len(jr) == 0 {
		return jr
	}

	sort.Sort(jr)

	pos := 0
	for i := range jr {
		if jr[pos].Name != jr[i].Name {
			pos++
		}
		jr[pos] = jr[i]
	}

	return jr[:pos+1]
}
