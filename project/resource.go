package project

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/concourse-friends/concourse-builder/model"
	"gopkg.in/yaml.v2"
)

type IInputResource interface {
	InputResources() JobResources
}

type Resource struct {
	// The name of the resource
	Name ResourceName

	// The type of the resource
	Type ResourceTypeName

	// The source of the resource
	Source IJobResourceSource

	// The scope of the resource
	Scope Scope

	// On what interval the resource to be pooled for updates
	CheckInterval model.Duration

	// Jobs needed to be part of the pipeline if this resource is consumed
	// Usually the job that produces the resource
	neededJobs Jobs
}

func (r *Resource) NeedJobs(jobs ...*Job) {
	r.neededJobs = append(r.neededJobs, jobs...)
}

func (r *Resource) NeededJobs() Jobs {
	registryType := GlobalTypeRegistry.RegisterType(r.Type)
	if registryType.Source == nil {
		return r.neededJobs
	}
	neededJobs := registryType.Source.NeededJobs()
	return append(neededJobs, r.neededJobs...)
}

type Resources []Resource

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
