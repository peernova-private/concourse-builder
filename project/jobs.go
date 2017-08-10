package project

import (
	"fmt"
	"sort"

	"github.com/concourse-friends/concourse-builder/model"
)

type Jobs []*Job

func (jobs Jobs) Len() int {
	return len(jobs)
}

func (jobs Jobs) Swap(i, j int) {
	jobs[i], jobs[j] = jobs[j], jobs[i]
}

func (jobs Jobs) Less(i, j int) bool {
	return jobs[i].Name < jobs[j].Name
}

func (jobs Jobs) SortByColumns() ([]Jobs, error) {
	blocked := make(map[*Job]struct{})
	mustHave := make(map[*Job]struct{})

	var columns []Jobs
	var column Jobs

	for _, job := range jobs {
		if len(job.AfterJobs) == 0 {
			column = append(column, job)
			continue
		}

		mustHave[job] = struct{}{}
		blocked[job] = struct{}{}
		for _, afterJob := range job.AfterJobs {
			if len(afterJob.AfterJobs) > 0 {
				blocked[afterJob] = struct{}{}
			}
		}
	}

	sort.Sort(column)
	columns = append(columns, column)

	for {
		var column Jobs
		move := false

		unblocked := make(map[*Job]struct{})
		for job := range blocked {
			stillBlocked := false
			for _, afterJob := range job.AfterJobs {
				if _, ok := blocked[afterJob]; ok {
					stillBlocked = true
					break
				}
			}
			if stillBlocked {
				continue
			}

			if _, ok := mustHave[job]; ok {
				column = append(column, job)
				delete(mustHave, job)
			}

			unblocked[job] = struct{}{}
			move = true
		}

		for job := range unblocked {
			delete(blocked, job)
		}

		if len(column) > 0 {
			sort.Sort(column)
			columns = append(columns, column)
		}

		if len(mustHave) == 0 {
			break
		}
		if !move {
			return nil, fmt.Errorf("There are jobs in circular positionaning")
		}
	}

	return columns, nil
}

func (jobs Jobs) NamesOfUsingResourceJobs(resource *JobResource) (model.JobNames, error) {
	var jobNames model.JobNames
	for _, job := range jobs {
		resources, err := job.Resources()
		if err != nil {
			return nil, err
		}
		found := false
		for _, jobResource := range resources {
			if resource.Name == jobResource.Name {
				found = true
				break
			}
		}
		if found {
			jobNames = append(jobNames, model.JobName(job.Name))
		}
	}

	return jobNames, nil
}
