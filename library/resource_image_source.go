package library

import "github.com/concourse-friends/concourse-builder/project"

type ResourceImageSource project.Resource

func (ris *ResourceImageSource) ModelSource() interface{} {
	return ris.Source.ModelSource()
}

func (ris *ResourceImageSource) NeededJobs() project.Jobs {
	return (*project.Resource)(ris).NeededJobs()
}
