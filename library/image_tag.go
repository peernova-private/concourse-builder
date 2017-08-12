package library

import "strings"

type ImageTag string

func ConvertToImageTag(raw string) ImageTag {
	tag := strings.Replace(raw, " ", "_", -1)
	return ImageTag(tag)
}
