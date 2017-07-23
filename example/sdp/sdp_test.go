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
jobs:
- name: git-image
  plan:
  - get: concourse-builder-git
  - put: git-image
    params:
      build: concourse-builder-git/template/sdp/docker/git-image
`

func TestSdp(t *testing.T) {
	prj, err := sdp.GenerateProject(&testSpecification{})
	assert.NoError(t, err)

	yml := &bytes.Buffer{}

	err = prj.Pipelines[0].Save(yml)
	assert.NoError(t, err)

	assert.Equal(t, expected, yml.String())
}
