package project

import "github.com/concourse-friends/concourse-builder/model"

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

func (job *Job) Model() (*model.Job, error) {
	var err error

	var modelSteps model.ISteps
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
