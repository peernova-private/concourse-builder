package sdpExample

import (
	"testing"

	"bytes"

	"github.com/ggeorgiev/concourse-builder/template/sdp"
	"github.com/stretchr/testify/assert"
)

var expected = `groups:
- name: images
  jobs:
  - git-image
`

func TestSdp(t *testing.T) {
	prj := sdp.GenerateProject(specification{})

	yml := &bytes.Buffer{}

	err := prj.Pipelines[0].Save(yml)
	assert.NoError(t, err)

	assert.Equal(t, expected, yml.String())
}
