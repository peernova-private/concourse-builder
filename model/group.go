package model

// A name of a group
type GroupName string

// A group of jobs to be visualized together
type Group struct {
	// The name of the group
	Name GroupName

	// Jobs that belong to the group
	Jobs []JobName
}

// Collection of groups
type Groups []*Group
