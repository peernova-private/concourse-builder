package library

import "path"

// Location is a struct that directs to file or directory inside a volume
type Location struct {
	Volume IVolume
	Path   string
}

func (l *Location) String() string {
	if l.Volume != nil {
		return path.Join(l.Volume.Path(), l.Path)
	}
	return l.Path
}
