package library

import "github.com/concourse-friends/concourse-builder/resource"

type S3PutParams struct {
	File string
}

func (s3pp *S3PutParams) ModelParams() interface{} {
	params := &resource.S3PutParams{
		File: s3pp.File,
	}

	return params
}
