package library

import (
	"strings"

	"github.com/concourse-friends/concourse-builder/library/primitive"
)

type ImageTag string

func ConvertToImageTag(raw string) ImageTag {
	tag := strings.Replace(raw, " ", "_", -1)
	tag = strings.Replace(tag, "/", "_", -1)
	return ImageTag(tag)
}

func BranchImageTag(branch *primitive.GitBranch) (ImageTag, bool) {
	if branch.IsMaster() || branch.IsImage() {
		return ConvertToImageTag(branch.FriendlyName()), true
	}
	return ConvertToImageTag(branch.BaseBranch()), false
}
