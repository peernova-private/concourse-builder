package project

import (
	"io"
	"sort"

	"github.com/concourse-friends/concourse-builder/model"
	yaml "gopkg.in/yaml.v2"
)

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

func (p *Pipeline) AllJobs() (Jobs, error) {
	checkJobs := make(JobsSet)

	for _, job := range p.Jobs {
		checkJobs[job] = struct{}{}
	}

	jobs := make(map[JobName]*Job)
	for job := checkJobs.Pop(); job != nil; job = checkJobs.Pop() {
		jobs[job.Name] = job

		resources, err := job.Resources()
		if err != nil {
			return nil, err
		}

		for _, resource := range resources {
			projResource := p.ResourceRegistry.MustGetResource(resource.Name)
			for _, resJob := range projResource.NeededJobs {
				if _, exists := jobs[resJob.Name]; exists {
					continue
				}
				checkJobs[resJob] = struct{}{}
				jobs[resJob.Name] = resJob
				job.AddJobToRunAfter(resJob)
			}
		}
	}

	sliceJobs := make(Jobs, 0, len(jobs))
	for _, job := range jobs {
		sliceJobs = append(sliceJobs, job)
	}

	return sliceJobs, nil
}

func (p *Pipeline) ModelGroups() (model.Groups, error) {
	groups := make(map[*JobGroup][]string)

	allJobs, err := p.AllJobs()
	if err != nil {
		return nil, err
	}

	allJobNames := make([]string, 0, len(allJobs))
	for _, job := range allJobs {
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

	groupsByOrder, err = SortJobGroups(groupsByOrder)
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
	allJobs, err := p.AllJobs()
	if err != nil {
		return nil, err
	}

	var jobResources JobResources
	for _, job := range allJobs {
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
	jobs, err := p.AllJobs()
	if err != nil {
		return nil, err
	}

	columns, err := jobs.SortByColumns()
	if err != nil {
		return nil, err
	}

	modelJobs := make(model.Jobs, 0, len(jobs))

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
