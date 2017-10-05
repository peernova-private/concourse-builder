package library

import (
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/resource"
)

type S3Source struct {
	// S3 Bucket
	Bucker *primitive.S3Bucket

	// Versioned file
	VersionedFile string
}

func (s3s *S3Source) ModelSource(scope project.Scope, info *project.ScopeInfo) interface{} {
	filePrefix := info.Scope(scope, "_")

	return &resource.S3Source{
		Bucket:          s3s.Bucker.Name,
		AccessKeyID:     s3s.Bucker.AccessKeyID,
		SecretAccessKey: s3s.Bucker.SecretAccessKey,
		RegionName:      s3s.Bucker.RegionName,
		VersionedFile:   filePrefix + s3s.VersionedFile,
	}
}
