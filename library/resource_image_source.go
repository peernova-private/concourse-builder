package library

import "github.com/concourse-friends/concourse-builder/project"

type ResourceImageSource project.Resource

func (ris *ResourceImageSource) ResourceScope() project.Scope {
	return project.Resource(*ris).Scope
}

func (ris *ResourceImageSource) ModelSource(scope project.Scope, info *project.ScopeInfo) interface{} {
	return ris.Source.ModelSource(scope, info)
}

func (ris *ResourceImageSource) NeededJobs() project.Jobs {
	return (*project.Resource)(ris).NeededJobs()
}
