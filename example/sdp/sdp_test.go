package sdpExample

import (
	"bytes"
	"testing"

	"github.com/concourse-friends/concourse-builder/resource"
	"github.com/concourse-friends/concourse-builder/template/sdp"
	"github.com/stretchr/testify/assert"
)

var expected = `groups:
- name: images
  jobs:
  - git-image
resources:
- name: git-image
  type: docker-image
  source: docker-io.peernova.com/concourse-builder/git-image
jobs:
- name: git-image
  plan:
  - put: git-image
    params:
      build: concourse-builder-git/template/sdp/docker/git-image
`

func TestSdp(t *testing.T) {
	prj := sdp.GenerateProject(&specification{
		ImageRepository: &resource.ImageRepository{
			Domain: "docker-io.peernova.com",
		},
	})

	yml := &bytes.Buffer{}

	err := prj.Pipelines[0].Save(yml)
	assert.NoError(t, err)

	assert.Equal(t, expected, yml.String())
}
