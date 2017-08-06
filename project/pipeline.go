package project

import (
	"io"
	"sort"

	"github.com/concourse-friends/concourse-builder/model"
	yaml "gopkg.in/yaml.v2"
)

type PipelineName string

type AllJobsGroupOption int

const (
	AllJobsGroupNone = iota
	AllJobsGroupFirst
	AllJobsGroupLast
)

type Pipeline struct {
	Name PipelineName

	AllJobsGroup AllJobsGroupOption
	Jobs         Jobs

	ResourceRegistry *ResourceRegistry
}

type Pipelines []*Pipeline

func NewPipeline() *Pipeline {
	return &Pipeline{
		ResourceRegistry: NewResourceRegistry(),
	}
}

func (p *Pipeline) ModelGroups() (model.Groups, error) {
	groups := make(map[*JobGroup][]string)

	allJobNames := make([]string, 0, len(p.Jobs))
	for _, job := range p.Jobs {
		allJobNames = append(allJobNames, string(job.Name))
		for _, group := range job.Groups {
			groups[group] = append(groups[group], string(job.Name))
		}
	}
	sort.Strings(allJobNames)

	groupsByOrder := JobGroups{}
	for group := range groups {
		groupsByOrder = append(groupsByOrder, group)
	}

	groupsByOrder, err := SortJobGroups(groupsByOrder)
	if err != nil {
		return nil, err
	}

	if p.AllJobsGroup != AllJobsGroupNone && len(groups) > 1 {
		all := &JobGroup{
			Name: "all",
		}

		if p.AllJobsGroup == AllJobsGroupFirst {
			groupsByOrder = append(JobGroups{all}, groupsByOrder...)
		} else /* if p.AllJobsGroup == AllJobsGroupLast */ {
			groupsByOrder = append(groupsByOrder, all)
		}
		groups[all] = allJobNames
	}

	modelGroups := model.Groups{}
	for _, group := range groupsByOrder {
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
	var jobResources JobResources
	for _, job := range p.Jobs {
		resources, err := job.Resources()
		if err != nil {
			return nil, err
		}
		jobResources = append(jobResources, resources...)
	}

	jobResources = jobResources.Deduplicate()

	var resources model.Resources
	for _, res := range jobResources {
		modelResource := res.Model(p.ResourceRegistry)
		resources = append(resources, modelResource)
	}

	return resources, nil
}

func (p *Pipeline) ModelJobs() (model.Jobs, error) {
	columns, err := p.Jobs.SortByColumns()
	if err != nil {
		return nil, err
	}

	modelJobs := make(model.Jobs, 0, len(p.Jobs))

	var previousColumn Jobs
	for _, column := range columns {
		for _, job := range column {
			modelJob, err := job.Model(previousColumn)
			if err != nil {
				return nil, err
			}
			modelJobs = append(modelJobs, modelJob)
		}
		previousColumn = column
	}

	return modelJobs, nil
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
