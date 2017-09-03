package primitive

import "regexp"

type GitBranch struct {
	Name string
}

var (
	branchPattern     = regexp.MustCompile(`^(.+?)(#(.+)|)$`)
	baseBranchPattern = regexp.MustCompile(`^.+#(.+)$`)
	isMasterPattern   = regexp.MustCompile(`^master$`)
	isReleasePattern  = regexp.MustCompile(`^release/.+`)
	isFeaturePattern  = regexp.MustCompile(`^feature/.+`)
	isTaskPattern     = regexp.MustCompile(`^task/.+`)
	isImagePattern    = regexp.MustCompile(`image`)
	isPrPattern       = regexp.MustCompile(`^(.+?)-pr(#(.+)|)$`)
)

func (gb *GitBranch) CanonicalName() string {
	return gb.Name
}

func (gb *GitBranch) FriendlyName() string {
	matches := branchPattern.FindAllStringSubmatch(gb.Name, -1)
	return matches[0][1]
}

func (gb *GitBranch) BaseBranch() string {
	matches := baseBranchPattern.FindAllStringSubmatch(gb.Name, -1)
	if matches != nil {
		return matches[0][1]
	}

	return "master"
}

func (gb *GitBranch) PrBranch() *GitBranch {
	name := gb.FriendlyName() + "-pr"
	base := gb.BaseBranch()
	if base != "" {
		name += "#" + base
	}

	return &GitBranch{
		Name: name,
	}
}

func (gb *GitBranch) IsMaster() bool {
	return isMasterPattern.MatchString(gb.Name)
}

func (gb *GitBranch) IsRelease() bool {
	return isReleasePattern.MatchString(gb.Name)
}

func (gb *GitBranch) IsFeature() bool {
	return isFeaturePattern.MatchString(gb.Name) && !isPrPattern.MatchString(gb.Name)
}

func (gb *GitBranch) IsTask() bool {
	return isTaskPattern.MatchString(gb.Name) && !isPrPattern.MatchString(gb.Name)
}

func (gb *GitBranch) IsImage() bool {
	return isImagePattern.MatchString(gb.Name)
}
