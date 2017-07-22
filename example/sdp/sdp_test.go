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
jobs:
- name: git-image
  plan:
  - put: git-image
`

func TestSdp(t *testing.T) {
	prj := sdp.GenerateProject(specification{})

	yml := &bytes.Buffer{}

	err := prj.Pipelines[0].Save(yml)
	assert.NoError(t, err)

	assert.Equal(t, expected, yml.String())
}
