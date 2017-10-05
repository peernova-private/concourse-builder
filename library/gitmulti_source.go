package library

import (
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type GitMultiSource struct {
	// Git repo and credentials
	Repo *primitive.GitRepo

	// Optional. Turns on multi-branch mode; takes a regular expression as argument --
	// branches matching the regular expression on origin will all be checked for changes.
	// Uses grep-style regular expression syntax
	Branches string
}

func (gms *GitMultiSource) ModelSource(scope project.Scope, info *project.ScopeInfo) interface{} {
	return &resource.GitMultibranchSource{
		URI:        gms.Repo.URI,
		PrivateKey: gms.Repo.PrivateKey,
		Branches:   gms.Branches,
	}
}
