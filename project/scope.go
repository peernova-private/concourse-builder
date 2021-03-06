package project

import "fmt"

type Scope int

type TeamName string
type InstallationName string

const (
	PipelineScope = iota
	AllPipelinesScope
	TeamScope
	AllTeamsScope
	InstallationScope
	UniverseScope
)

type ScopeInfo struct {
	Pipeline     PipelineName
	Team         TeamName
	Installation InstallationName
}

func (info *ScopeInfo) Scope(scope Scope, delimiter string) string {
	switch scope {
	case UniverseScope:
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
	case AllTeamsScope:
		return fmt.Sprintf("%s%s",
			info.Installation, delimiter)
	case AllPipelinesScope:
		return fmt.Sprintf("%s%s%s%s",
			info.Installation, delimiter,
			info.Team, delimiter)
	}

	return ""
}
