package model

import (
	"strings"
	"time"
)

type Duration time.Duration

func (d Duration) MarshalYAML() (interface{}, error) {
	str := time.Duration(d).String()
	if strings.HasSuffix(str, "h0m0s") {
		return str[:len(str)-4], nil
	}
	if strings.HasSuffix(str, "m0s") {
		return str[:len(str)-2], nil
	}
	return str, nil
}
