package project

type JobName string

type JobGroup struct {
	Name     string
	Position int
}

type JobGroups []JobGroup

func (a JobGroups) Len() int           { return len(a) }
func (a JobGroups) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a JobGroups) Less(i, j int) bool { return a[i].Position < a[j].Position }

type Job struct {
	Name   JobName
	Groups []JobGroup
}
