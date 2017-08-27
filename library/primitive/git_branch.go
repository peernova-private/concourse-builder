package primitive

import "regexp"

type GitBranch struct {
	Branch string
}

func (gb *GitBranch) Name() string {
	return gb.Branch
}

var isMasterPattern = regexp.MustCompile(`^master$`)

func (gb *GitBranch) IsMaster() bool {
	return isMasterPattern.Match([]byte(gb.Branch))
}

var isReleasePattern = regexp.MustCompile(`^release/.+`)

func (gb *GitBranch) IsRelease() bool {
	return isReleasePattern.Match([]byte(gb.Branch))
}

var isFeaturePattern = regexp.MustCompile(`^feature/.+`)

func (gb *GitBranch) IsFeature() bool {
	return isFeaturePattern.Match([]byte(gb.Branch))
}

var isTaskPattern = regexp.MustCompile(`^task/.+`)

func (gb *GitBranch) IsTask() bool {
	return isTaskPattern.Match([]byte(gb.Branch))
}
