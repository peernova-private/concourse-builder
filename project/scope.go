package project

import "fmt"

type Scope int

type TeamName string
type InstallationName string

const (
	PipelineScope = iota
	TeamScope
	InstallationScope
	UniversalScope
)

type ScopeInfo struct {
	Pipeline     PipelineName
	Team         TeamName
	Installation InstallationName
}

func (info *ScopeInfo) Scope(scope Scope, delimiter string) string {
	switch scope {
	case UniversalScope:
		return fmt.Sprintf("%s%s%s%s%s%s",
			info.Installation, delimiter,
			info.Team, delimiter,
			info.Pipeline, delimiter)
	case InstallationScope:
		return fmt.Sprintf("%s%s%s%s",
			info.Team, delimiter,
			info.Pipeline, delimiter)
	case TeamScope:
		return fmt.Sprintf("%s%s",
			info.Pipeline, delimiter)
	}

	return ""
}
