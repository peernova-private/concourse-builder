package primitive

import (
	"path"

	"github.com/concourse-friends/concourse-builder/project"
)

// Location is a struct that directs to file or directory inside a volume
type Location struct {
	Volume       IVolume
	RelativePath string
}

func (l *Location) Path() string {
	if l.Volume != nil {
		return path.Join(l.Volume.Path(), l.RelativePath)
	}
	return l.RelativePath
}

func (l *Location) InputResources() project.JobResources {
	var resources project.JobResources
	if res, ok := l.Volume.(*project.JobResource); ok {
		resources = append(resources, res)
	}
	return resources
}

func (l *Location) OutputNames() []string {
	if output, ok := l.Volume.(*project.TaskOutput); ok {
		return []string{output.Name()}
	}
	if res, ok := l.Volume.(*project.JobResource); ok {
		return []string{string(res.Name)}
	}
	return nil
}

func (l *Location) Value() string {
	return l.Path()
}
