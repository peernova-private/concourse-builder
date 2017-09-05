package project

import "github.com/concourse-friends/concourse-builder/model"

type IStep interface {
	Model() (model.IStep, error)
	InputResources() (JobResources, error)
	OutputResource() (*Resource, error)
}

type ISteps []IStep
