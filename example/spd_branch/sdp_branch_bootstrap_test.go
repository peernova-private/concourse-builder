package sdpBranchExample

import (
	"bytes"
	"testing"

	"github.com/concourse-friends/concourse-builder/template/sdp_branch"
	"github.com/concourse-friends/concourse-builder/test"
	"github.com/stretchr/testify/assert"
)

var buildImageScript = `
          H4sIAAAAAAAC/4STb2vbMBDGX1ef4qkXKAw8Mxh7l0JI3C60JcXp9mYbqSKda1HPEpKS/iEfftgm8Z94\
          2Uudnvvd3SPdh/No42y0VkVExRZr7jKWLBYP40fzIh8Zc+QR0itj02/x9GY1Sa6XY283xJhK8RPhO4LR\
          bDG9iZOr+W28ms2TACE/ii8f4vtlgN/MZ1SwMxKZRnCfE3cEZ0io9A1dDLRFH4CyxS23iq9zCvaYXl6N\
          U+TgM4JUloTX9g0vGVmqY1o8k01VTnCejAO3g7S6pnLgpS/09Que3pWB81YVT9DpICxgZy2rUp47Yqlq\
          23WVLO5W87vJdfwfQxrh8OBz3xvWktFOVdN6jTVh40hCFdVlSYPI+cbRyR6DUXMXYDxGUAlavb4qj89V\
          wp9nqSxCA2PJcEuyRpwP/4s9QBj07qKPDaHkptpiufieTOsnVQVGzXHJpGYAcFQ9aqkqhTAIk3Zu9GlY\
          LHVBjAmJpg9WmVy51noyXGJ2ePFactkJdRyIf0xuD3OXDdGW54d4N7Gce4B3ytBy0Qijcl+jvqP1z26V\
          FtyflB5384/S3V2upqq+44BmB28RSlz8Ki6w229RGEoSWhJ29UKFIpXH5f8GAAD//2JqijKaBAAA |\`

var expectedBootstrap = `groups:
- name: all
  jobs:
  - curl-image
  - fly-image
  - self-update
- name: images
  jobs:
  - curl-image
  - fly-image
- name: sys
  jobs:
  - curl-image
  - fly-image
resource_types:
- name: dummy
  type: docker-image
  source:
    repository: registry.com/concourse-builder/dummy_resource-image
    tag: master
    aws_access_key_id: key
    aws_secret_access_key: secret
resources:
- name: concourse-builder-git
  type: git
  source:
    uri: git@github.com:concourse-friends/concourse-builder.git
    branch: master
    private_key: private-key
- name: curl-image
  type: docker-image
  source:
    repository: registry.com/concourse-builder/curl-image
    tag: master
    aws_access_key_id: key
    aws_secret_access_key: secret
- name: fly-image
  type: docker-image
  source:
    repository: registry.com/concourse-builder/fly-image
    tag: master
    aws_access_key_id: key
    aws_secret_access_key: secret
- name: go-image
  type: docker-image
  source:
    repository: golang
  check_every: 24h
- name: pipeline
  type: dummy
- name: ubuntu-image
  type: docker-image
  source:
    repository: ubuntu
  check_every: 24h
jobs:
- name: curl-image
  plan:
  - get: ubuntu-image
    trigger: true
  - task: prepare
    image: ubuntu-image
    config:
      platform: linux
      params:
        DOCKERFILE_STEPS: |-
          H4sIAAAAAAAC/1TMsQrCUAyF4b1PcUDoIIS8iYPg1iXWKBdiCUmu6NuLlg5dP/5zzpcTUgukb0zDAW3J
          EjPMPQzTMI4QL3pooftNSve21fTZBr+P2VSW7vv0jyvFExR38EuCrV1ZvNhaVvLxGwAA//9QRiyQjwAA
          AA==
        FROM_IMAGE: ubuntu
      run:
        path: /bin/bash
        args:
        - -c
        - |-
          mkdir -p /tmp \
          && echo \` + buildImageScript + `
              base64 --decode |\
              gzip -cfd > /tmp/script.sh \
          && cat /tmp/script.sh \
          && echo \
          && chmod 755 /tmp/script.sh \
          && /tmp/script.sh
      outputs:
      - name: prepared
        path: prepared
  - put: curl-image
    params:
      build: prepared
    get_params:
      skip_download: true
- name: fly-image
  plan:
  - aggregate:
    - get: concourse-builder-git
      trigger: true
    - get: curl-image
      trigger: true
      passed:
      - curl-image
      params:
        save: true
  - task: prepare
    image: curl-image
    config:
      platform: linux
      inputs:
      - name: concourse-builder-git
      params:
        DOCKERFILE_DIR: concourse-builder-git/docker/fly
        EVAL: echo ENV FLY_VERSION=` + "`" + `curl http://concourse.com/api/v1/info | awk -F
          ',' ' { print $1 } ' | awk -F ':' ' { print $2 } '` + "`" + `
        FROM_IMAGE: registry.com/concourse-builder/curl-image:master
      run:
        path: /bin/bash
        args:
        - -c
        - |-
          mkdir -p /tmp \
          && echo \` + buildImageScript + `
              base64 --decode |\
              gzip -cfd > /tmp/script.sh \
          && cat /tmp/script.sh \
          && echo \
          && chmod 755 /tmp/script.sh \
          && /tmp/script.sh
      outputs:
      - name: prepared
        path: prepared
  - put: fly-image
    params:
      build: prepared
      load_base: curl-image
    get_params:
      skip_download: true
- name: self-update
  plan:
  - aggregate:
    - get: concourse-builder-git
      trigger: true
      passed:
      - fly-image
    - get: fly-image
      trigger: true
      passed:
      - fly-image
    - get: go-image
      trigger: true
  - task: check
    image: fly-image
    config:
      platform: linux
      params:
        CONCOURSE_URL: http://concourse.com
      run:
        path: /bin/fly/check_version.sh
  - task: prepare pipelines
    image: go-image
    config:
      platform: linux
      inputs:
      - name: concourse-builder-git
      params:
        BRANCH: branch
        PIPELINES: pipelines
      run:
        path: concourse-builder-git/foo
      outputs:
      - name: pipelines
        path: pipelines
  - task: update pipelines
    image: fly-image
    config:
      platform: linux
      inputs:
      - name: pipelines
      params:
        CONCOURSE_PASSWORD: password
        CONCOURSE_URL: http://concourse.com
        CONCOURSE_USER: user
        PIPELINES: pipelines
      run:
        path: /bin/fly/set_pipelines.sh
  - put: pipeline
`

func TestSdpBranchBootstrap(t *testing.T) {
	prj, err := sdpBranch.GenerateBootstrapProject(&testSpecification{})
	assert.NoError(t, err)

	yml := &bytes.Buffer{}

	err = prj.Pipelines[0].Save(yml)
	assert.NoError(t, err)

	test.AssertEqual(t, expectedBootstrap, yml.String())
}
