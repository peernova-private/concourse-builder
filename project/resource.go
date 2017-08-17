package project

import (
	"fmt"
	"strings"

	"crypto/sha256"

	"github.com/concourse-friends/concourse-builder/model"
	"gopkg.in/yaml.v2"
	"encoding/hex"
)

type Resource struct {
	// The name of the resource
	Name ResourceName

	// The type of the resource
	Type model.ResourceTypeName

	// The sourse of the resource
	Source IJobResourceSource

	// On what interval the resource to be pooled for updates
	CheckInterval model.Duration

	// Jobs needed to be part of the pipeline if this resource is consumed
	// Usually the job that produces the resource
	NeededJobs Jobs
}

type ResourceHash string

func (r *Resource) MustHash() ResourceHash {
	yml, err := yaml.Marshal(r)
	if err != nil {
		panic(err.Error())
	}

	str := strings.Replace(string(yml), fmt.Sprintf("name: %s", r.Name), "@@@", 1)

	sha := sha256.New()
	sha.Write([]byte(str))
	hash := sha.Sum(nil)

	return ResourceHash(hex.EncodeToString(hash))
}
