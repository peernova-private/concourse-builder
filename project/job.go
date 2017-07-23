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

func (job *Job) InputResources() (model.Resources, error) {
	var resources model.Resources

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

	return resources, nil
}

func (job *Job) Resources() (model.Resources, error) {
	var resources model.Resources

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
		resources = append(resources, outputResource)
	}

	return resources, nil
}

func (job *Job) Model() (*model.Job, error) {
	var err error

	var modelSteps model.ISteps

	inputs, err := job.InputResources()
	if err != nil {
		return nil, err
	}

	for _, input := range inputs {
		step := &model.Get{
			Get: input.Name,
		}
		modelSteps = append(modelSteps, step)
	}

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
