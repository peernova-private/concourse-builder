package library

import (
	"github.com/concourse-friends/concourse-builder/resource"
)

type GitMultiSource struct {
	// Git repo and credentials
	Repo *GitRepo

	// Optional. Turns on multi-branch mode; takes a regular expression as argument --
	// branches matching the regular expression on origin will all be checked for changes.
	// Uses grep-style regular expression syntax
	Branches string
}

func (gms *GitMultiSource) ModelSource() interface{} {
	return &resource.GitMultibranchSource{
		URI:        gms.Repo.URI,
		PrivateKey: gms.Repo.PrivateKey,
		Branches:   gms.Branches,
	}
}
