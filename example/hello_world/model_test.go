package hello_world

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var yml = `jobs:
- name: hello-world
  plan:
  - task: say-hello
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: ubuntu
      run:
        path: echo
        args:
        - Hello, world!
`

func TestHelloWorld(t *testing.T) {
	result, err := yaml.Marshal(Pipeline)
	assert.Nil(t, err)

	assert.Equal(t, yml, string(result))
}
