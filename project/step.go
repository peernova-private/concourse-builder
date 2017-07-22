package project

import "github.com/concourse-friends/concourse-builder/model"

type IStep interface {
	Model() (model.IStep, error)
}

type ISteps []IStep
