package project

import "github.com/concourse-friends/concourse-builder/model"

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
