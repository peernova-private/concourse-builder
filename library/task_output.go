package library

import (
	"strings"
)

type TaskOutput struct {
	Directory string
}

func (to *TaskOutput) Name() string {
	return strings.Split(to.Directory, "/")[0]
}

func (to *TaskOutput) Path() string {
	return to.Directory
}
