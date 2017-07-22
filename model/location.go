package model

import "path"

// Location is a struct that directs to file or directory inside a volume
type Location struct {
	Volume IVolume
	Path   string
}

func (l *Location) MarshalYAML() (interface{}, error) {
	return path.Join(l.Volume.Path(), l.Path), nil
}
