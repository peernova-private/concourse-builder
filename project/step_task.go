package project

import (
	"fmt"
	"sort"

	"github.com/concourse-friends/concourse-builder/model"
)

type IRun interface {
	Path() string
}

type IOutput interface {
	Name() string
	Path() string
}

type ITaskInput interface {
	OutputNames() []string
}

type IEnvironmentValue interface {
	Value() string
}

type ITaskDirectory interface {
	Path() string
}

type TaskStep struct {
	Platform    model.Platform
	Name        model.TaskName
	Image       *JobResource
	Privileged  bool
	Run         IRun
	Outputs     []IOutput
	Environment map[string]interface{}
	Directory   ITaskDirectory
	Arguments   []interface{}
}

func (ts *TaskStep) Model() (model.IStep, error) {
	task := &model.Task{
		Task:       ts.Name,
		Privileged: ts.Privileged,
		Config: &model.TaskConfig{
			Platform: ts.Platform,
			Run: &model.TaskRun{
				Path: ts.Run.Path(),
			},
			Params: make(map[string]interface{}),
		},
	}

	if ts.Directory != nil {
		task.Config.Run.Dir = ts.Directory.Path()
	}

	for _, argument := range ts.Arguments {
		switch value := argument.(type) {
		case string:
			task.Config.Run.Args = append(task.Config.Run.Args, value)
		default:
			panic(fmt.Sprintf("%v is unsuported type", argument))
		}
	}

	if ts.Image != nil {
		task.Image = model.ResourceName(ts.Image.Name)
	}

	inputResources, err := ts.ExecutionResources()
	if err != nil {
		return nil, err
	}

	inputsMap := make(map[string]struct{})
	for _, inputResource := range inputResources {
		inputsMap[string(inputResource.Name)] = struct{}{}
	}

	if directory, ok := ts.Directory.(ITaskInput); ok {
		names := directory.OutputNames()
		for _, name := range names {
			inputsMap[name] = struct{}{}
		}
	}

	// TODO: revisit this one
	for _, value := range ts.Environment {
		var names []string
		if variable, ok := value.(ITaskInput); ok {
			names = variable.OutputNames()
		}
		for _, name := range names {
			inputsMap[name] = struct{}{}
		}
	}

	inputs := make([]string, 0, len(inputsMap))
	for input := range inputsMap {
		inputs = append(inputs, input)
	}

	sort.Strings(inputs)
	for _, name := range inputs {
		task.Config.Inputs = append(task.Config.Inputs, &model.TaskInput{
			Name: model.ResourceName(name),
		})
	}

	for _, output := range ts.Outputs {
		task.Config.Outputs = append(task.Config.Outputs, &model.TaskOutput{
			Name: output.Name(),
			Path: output.Path(),
		})
	}

	for name, value := range ts.Environment {
		if param, ok := value.(IEnvironmentValue); ok {
			task.Config.Params[name] = param.Value()
		} else {
			task.Config.Params[name] = value
		}
	}

	return task, nil
}

func (ts *TaskStep) ExecutionResources() (JobResources, error) {
	var resources JobResources

	if locationResources, ok := ts.Run.(IInputResource); ok {
		resources = append(resources, locationResources.InputResources()...)
	}

	for _, value := range ts.Environment {
		if variable, ok := value.(IInputResource); ok {
			resources = append(resources, variable.InputResources()...)
		}
	}

	if directory, ok := ts.Directory.(IInputResource); ok {
		resources = append(resources, directory.InputResources()...)
	}

	for _, argument := range ts.Arguments {
		if argResource, ok := argument.(IInputResource); ok {
			resources = append(resources, argResource.InputResources()...)
		}
	}

	return resources.Deduplicate(), nil
}

func (ts *TaskStep) InputResources() (JobResources, error) {
	var resources JobResources

	if ts.Image != nil {
		resources = append(resources, ts.Image)
	}

	executionResources, err := ts.ExecutionResources()
	if err != nil {
		return nil, err
	}

	resources = append(resources, executionResources...)
	return resources.Deduplicate(), nil
}

func (ts *TaskStep) OutputResource() (*JobResource, error) {
	return nil, nil
}
