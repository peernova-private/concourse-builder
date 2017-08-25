package project

import (
	"fmt"
	"sort"

	"github.com/concourse-friends/concourse-builder/model"
)

type IRun interface {
	InputResources() JobResources
	Path() string
}

type IOutput interface {
	Name() string
	Path() string
}

type IEnvironmentValue interface {
	Value() interface{}
}

type IEnvironmentInput interface {
	OutputName() string
}

type IEnvironmentResource interface {
	InputResources() JobResources
}

type ITaskDirectory interface {
	Path() string
}

type ITaskDirectoryResource interface {
	InputResources() JobResources
}

type IArgumentResource interface {
	InputResources() JobResources
}

type TaskStep struct {
	Platform    model.Platform
	Name        model.TaskName
	Image       *JobResource
	Run         IRun
	Outputs     []IOutput
	Environment map[string]interface{}
	Directory   ITaskDirectory
	Arguments   []interface{}
}

func (ts *TaskStep) Model() (model.IStep, error) {
	task := &model.Task{
		Task: ts.Name,
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

	// TODO: revisit this one
	for _, value := range ts.Environment {
		if param, ok := value.(IEnvironmentInput); ok {
			name := param.OutputName()
			if name != "" {
				inputsMap[name] = struct{}{}
			}
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

	locationResources := ts.Run.InputResources()
	resources = append(resources, locationResources...)

	for _, value := range ts.Environment {
		if variable, ok := value.(IEnvironmentResource); ok {
			resources = append(resources, variable.InputResources()...)
		}
	}

	if directory, ok := ts.Directory.(ITaskDirectoryResource); ok {
		resources = append(resources, directory.InputResources()...)
	}

	for _, argument := range ts.Arguments {
		if argResource, ok := argument.(IArgumentResource); ok {
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
