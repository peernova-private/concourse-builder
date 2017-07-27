package sdpExample

import (
	"bytes"
	"testing"

	"github.com/concourse-friends/concourse-builder/template/sdp"
	"github.com/stretchr/testify/assert"
)

var expected = `groups:
- name: images
  jobs:
  - git-image
resources:
- name: concourse-builder-git
  type: git
  source:
    uri: git@github.com:concourse-friends/concourse-builder.git
    branch: master
    private_key: private-key
- name: git-image
  type: docker-image
  source:
    repository: repository.com/concourse-builder/git-image
- name: ubuntu-image
  type: docker-image
  source:
    repository: ubuntu
    tag: "16.04"
  check_every: 24h
jobs:
- name: git-image
  plan:
  - aggregate:
    - get: concourse-builder-git
      trigger: true
    - get: ubuntu-image
      trigger: true
  - task: prepare
    image: ubuntu-image
    config:
      platform: linux
      inputs:
      - name: concourse-builder-git
      params:
        DOCKER_STEPS: concourse-builder-git/docker/git_steps
        FROM_IMAGE: ubuntu:16.04
      run:
        path: concourse-builder-git/scripts/docker_image_prepare.sh
      outputs:
      - name: prepared
        path: prepared
  - put: git-image
    params:
      build: prepared
    get_params:
      skip_download: true
`

func TestSdp(t *testing.T) {
	prj, err := sdp.GenerateProject(&testSpecification{})
	assert.NoError(t, err)

	yml := &bytes.Buffer{}

	err = prj.Pipelines[0].Save(yml)
	assert.NoError(t, err)

	assert.Equal(t, expected, yml.String())
}
