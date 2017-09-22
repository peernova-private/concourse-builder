package project

import (
	"log"

	"github.com/concourse-friends/concourse-builder/model"
)

type JobName string

type Job struct {
	Name           JobName
	Groups         JobGroups
	Sequentiality  Sequentiality
	ExtraResources JobResources
	Steps          ISteps
	OnSuccess      IStep
	OnFailure      IStep
	AfterJobs      map[*Job]struct{}
}

func (job *Job) AddToGroup(groups ...*JobGroup) {
	job.Groups = append(job.Groups, groups...)
}

func (job *Job) AddJobToRunAfter(jobs ...*Job) {
	if job.AfterJobs == nil {
		job.AfterJobs = make(map[*Job]struct{})
	}

	for _, afterJob := range jobs {
		if _, ok := job.AfterJobs[afterJob]; !ok {
			log.Printf("Run job %s after %s", job.Name, afterJob.Name)
			job.AfterJobs[afterJob] = struct{}{}
		}
	}
}

func (job *Job) InputResources() (JobResources, error) {
	var resources JobResources
	resources = append(resources, job.ExtraResources...)

	steps := append(ISteps{job.OnSuccess, job.OnFailure}, job.Steps...)
	for _, step := range steps {
		if step == nil {
			continue
		}
		inputResources, err := step.InputResources()
		if err != nil {
			return nil, err
		}
		resources = append(resources, inputResources...)
	}

	return resources.Deduplicate(), nil
}

func (job *Job) Resources() (JobResources, error) {
	var resources JobResources
	resources = append(resources, job.ExtraResources...)

	steps := append(ISteps{job.OnSuccess, job.OnFailure}, job.Steps...)
	for _, step := range steps {
		if step == nil {
			continue
		}
		inputResources, err := step.InputResources()
		if err != nil {
			return nil, err
		}
		resources = append(resources, inputResources...)

		outputResource, err := step.OutputResource()
		if err != nil {
			return nil, err
		}
		if outputResource != nil {
			resources = append(resources, &JobResource{Name: outputResource.Name})
		}
	}

	return resources.Deduplicate(), nil
}

func (job *Job) Model(previousColumns []Jobs) (*model.Job, error) {
	var err error

	var modelSteps model.ISteps

	inputs, err := job.InputResources()
	if err != nil {
		return nil, err
	}

	var modelGetSteps model.ISteps
	for _, input := range inputs {
		step := &model.Get{
			Get:     model.ResourceName(input.Name),
			Trigger: input.Trigger,
			Params:  input.GetParams,
		}

		for c := len(previousColumns) - 1; c >= 0; c-- {
			step.Passed, err = previousColumns[c].NamesOfUsingResourceJobs(input)
			if err != nil {
				return nil, err
			}

			if step.Passed != nil {
				break
			}
		}

		modelGetSteps = append(modelGetSteps, step)
	}

	if len(modelGetSteps) > 1 {
		modelGetSteps = model.ISteps{
			&model.Aggregation{
				Aggregate: modelGetSteps,
			},
		}
	}
	modelSteps = append(modelSteps, modelGetSteps...)

	for _, step := range job.Steps {
		modelStep, err := step.Model()
		if err != nil {
			return nil, err
		}
		modelSteps = append(modelSteps, modelStep)
	}

	var modelOnSuccessStep model.IStep
	if job.OnSuccess != nil {
		modelOnSuccessStep, err = job.OnSuccess.Model()
		if err != nil {
			return nil, err
		}
	}

	var modelOnFailureStep model.IStep
	if job.OnFailure != nil {
		modelOnFailureStep, err = job.OnFailure.Model()
		if err != nil {
			return nil, err
		}
	}

	return &model.Job{
		Name:      model.JobName(job.Name),
		Plan:      modelSteps,
		OnSuccess: modelOnSuccessStep,
		OnFailure: modelOnFailureStep,
	}, nil
}

func (job *Job) TaskOutputIndex(output *TaskOutput) (int, int) {
	for s := range job.Steps {
		task, ok := job.Steps[s].(*TaskStep)
		if !ok {
			continue
		}

		for o := range task.Outputs {
			out, ok := task.Outputs[o].(*TaskOutput)
			if !ok {
				continue
			}

			if out == output {
				return s, o
			}
		}
	}

	return 0, 0
}
