package project

import (
	"io"

	"sort"

	"github.com/ggeorgiev/concourse-builder/model"
	yaml "gopkg.in/yaml.v2"
)

type PipelineName string

type Pipeline struct {
	Name PipelineName
	Jobs []*Job
}

func (p *Pipeline) ModelGroups() (model.Groups, error) {
	groups := make(map[JobGroup][]string)

	for _, job := range p.Jobs {
		for _, group := range job.Groups {
			groups[group] = append(groups[group], string(job.Name))
		}
	}

	groupsNames := JobGroups{}
	for group := range groups {
		groupsNames = append(groupsNames, group)
	}
	sort.Sort(groupsNames)

	modelGroups := model.Groups{}
	for _, group := range groupsNames {
		modelJobs := model.JobNames{}

		jobNames := groups[group]
		sort.Strings(jobNames)
		for _, jobName := range jobNames {
			modelJobs = append(modelJobs, model.JobName(jobName))
		}

		modelGroup := model.Group{
			Name: model.GroupName(group.Name),
			Jobs: modelJobs,
		}

		modelGroups = append(modelGroups, &modelGroup)
	}

	return modelGroups, nil
}

func (p *Pipeline) ModelResourceTypes() (model.ResourceTypes, error) {
	return nil, nil
}

func (p *Pipeline) ModelResources() (model.Resources, error) {
	return nil, nil
}

func (p *Pipeline) ModelJobs() (model.Jobs, error) {
	return nil, nil
}

func (p *Pipeline) Save(writer io.Writer) error {
	groups, err := p.ModelGroups()
	if err != nil {
		return err
	}

	resourceTypes, err := p.ModelResourceTypes()
	if err != nil {
		return err
	}

	resources, err := p.ModelResources()
	if err != nil {
		return err
	}

	jobs, err := p.ModelJobs()
	if err != nil {
		return err
	}

	pipeline := model.Pipeline{
		Groups:        groups,
		ResourceTypes: resourceTypes,
		Resources:     resources,
		Jobs:          jobs,
	}

	result, err := yaml.Marshal(pipeline)
	if err != nil {
		return err
	}

	_, err = writer.Write(result)
	if err != nil {
		return err
	}

	return nil
}
