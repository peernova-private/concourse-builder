package library

// Interface of object that has location in the running container.
// This include but not is not limited to: Task outputs, resources and preinstalled docker places.
type IVolume interface {
	Path() string
}
