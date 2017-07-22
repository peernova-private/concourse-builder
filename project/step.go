package project

import "github.com/concourse-friends/concourse-builder/model"

type IStep interface {
	Model() (model.IStep, error)
	InputResources() (model.Resources, error)
	OutputResource() (*model.Resource, error)
}

type ISteps []IStep
