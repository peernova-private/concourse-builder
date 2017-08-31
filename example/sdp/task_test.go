// Copyright (C) 2017 PeerNova, Inc.
//
// All rights reserved.
//
// PeerNova and Cuneiform are trademarks of PeerNova, Inc. References to
// third-party marks or brands are the property of their respective owners.
// No rights or licenses are granted, express or implied, unless set forth in
// a written agreement signed by PeerNova, Inc. You may not distribute,
// disseminate, copy, record, modify, enhance, supplement, create derivative
// works from, adapt, or translate any content contained herein except as
// otherwise expressly permitted pursuant to a written agreement signed by
// PeerNova, Inc.

package sdpExample

import (
	"testing"
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var expectedTask = `task: test_task
timeout: 10s
image: image_resource
config:
  platform: linux
  inputs:
  - name: concourse_git
    path: concourse_git
  - name: go_ci_tool_image
    path: go_ci_tool_image
  params:
    param_1: 1
    param_2: "2"
  run:
    path: script_path
    args:
    - a
    - b
    - c
    dir: script_dir
  outputs:
  - name: go_artifacts
    path: go_artifacts
ensure:
  put: concourse_master_pr
  attempts: 2
  timeout: 10s
  params:
    param_1: 1
    param_2: "2"
on_success:
  put: concourse_master_pr
  attempts: 2
  timeout: 10s
  params:
    param_1: 1
    param_2: "2"
on_failure:
  put: concourse_master_pr
  attempts: 3
  timeout: 10s
  params:
    context: '@code_gen'
    status: failure
  get_params:
    git.disable_lfs: true
attempts: 1
`

func TestConstructTask(t *testing.T) {
	simpleTask := &model.Task{
		Task:      "test_task",
		Timeout:   10 * time.Second,
		Image:     "image_resource",
		Config:    prepareTaskConfig(),
		Attempts:  1,
		OnSuccess: prepareOnSuccessStep(),
		OnFailure: prepareOnFailureStep(),
		Ensure:    prepareEnsureStep(),
	}

	taskString, err := yaml.Marshal(simpleTask)
	assert.NoError(t, err)
	test.AssertEqual(t, expectedTask, string(taskString))

}

func prepareTaskConfig() *model.TaskConfig {
	params := map[string]interface{}{
		"param_1": 1,
		"param_2": "2",
	}

	tc := &model.TaskConfig{
		Platform: "linux",
		Params:   params,
		Inputs:   prepareTaskInput(),
		Run:      prepareTaskRun(),
		Outputs:  prepareTaskOutput(),
	}

	return tc
}

func prepareTaskInput() []*model.TaskInput {
	inputs := []string{"concourse_git", "go_ci_tool_image"}
	taskInputs := make([]*model.TaskInput, 2)

	for i, v := range inputs {
		taskInputs[i] = &model.TaskInput{
			Name: model.ResourceName(v),
			Path: v,
		}
	}

	return taskInputs
}

func prepareTaskRun() *model.TaskRun {
	tr := &model.TaskRun{
		Path: "script_path",
		Args: []string{"a", "b", "c"},
		Dir:  "script_dir",
	}

	return tr
}

func prepareTaskOutput() []*model.TaskOutput {
	inputs := []string{"go_artifacts"}
	taskInputs := make([]*model.TaskOutput, 1)

	for i, v := range inputs {
		taskInputs[i] = &model.TaskOutput{
			Name: v,
			Path: v,
		}
	}

	return taskInputs
}

func prepareOnSuccessStep() model.IStep {
	params := map[string]interface{}{
		"param_1": 1,
		"param_2": "2",
	}

	putStep := &model.Put{
		Put:      "concourse_master_pr",
		Attempts: 2,
		Params:   params,
		Timeout:  10 * time.Second,
	}

	return putStep
}

func prepareOnFailureStep() model.IStep {
	params := map[string]interface{}{
		"context": "@code_gen",
		"status":  "failure",
	}

	getParams := map[string]interface{}{
		"git.disable_lfs": true,
	}

	putStep := &model.Put{
		Put:       "concourse_master_pr",
		Attempts:  3,
		Params:    params,
		Timeout:   10 * time.Second,
		GetParams: getParams,
	}

	return putStep
}

func prepareEnsureStep() model.IStep {
	params := map[string]interface{}{
		"param_1": 1,
		"param_2": "2",
	}

	putStep := &model.Put{
		Put:      "concourse_master_pr",
		Attempts: 2,
		Params:   params,
		Timeout:  10 * time.Second,
	}

	return putStep
}
