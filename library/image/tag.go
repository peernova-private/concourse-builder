package image

import (
	"strings"

	"github.com/concourse-friends/concourse-builder/library/primitive"
)

type Tag string

func ConvertToImageTag(raw string) Tag {
	tag := strings.Replace(raw, " ", "_", -1)
	tag = strings.Replace(tag, "/", "_", -1)
	return Tag(tag)
}

func BranchImageTag(branch *primitive.GitBranch) (Tag, bool) {
	if branch.IsMaster() || branch.IsImage() {
		return ConvertToImageTag(branch.FriendlyName()), true
	}
	return ConvertToImageTag(branch.BaseBranch()), false
}
