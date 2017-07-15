package project

type JobName string

type Job struct {
	Name   JobName
	Groups JobGroups
}

type Jobs []*Job
