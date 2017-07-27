package project

import (
	"github.com/concourse-friends/concourse-builder/model"
)

type JobName string

type Job struct {
	Name          JobName
	Groups        JobGroups
	Sequentiality Sequentiality
	Steps         ISteps
	OnSuccess     IStep
	OnFailure     IStep
}

type Jobs []*Job

func (job *Job) InputResources() (JobResources, error) {
	var resources JobResources

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
			resources = append(resources, outputResource)
		}
	}

	return resources.Deduplicate(), nil
}

func (job *Job) Model() (*model.Job, error) {
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
