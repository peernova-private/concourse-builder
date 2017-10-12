package image

import (
	"strings"
)

type Tag string

func ConvertToImageTag(raw string) Tag {
	tag := strings.Replace(raw, " ", "_", -1)
	tag = strings.Replace(tag, "/", "_", -1)
	return Tag(tag)
}
