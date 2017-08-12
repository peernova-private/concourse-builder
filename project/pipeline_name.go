package project

import "strings"

type PipelineName string

func ConvertToPipelineName(raw string) PipelineName {
	name := strings.Replace(raw, " ", "_", -1)
	name = strings.Replace(name, "/", "_", -1)
	return PipelineName(name)
}
