package project

import "fmt"

type JobGroup struct {
	Name  string
	After *JobGroup
}

type JobGroups []*JobGroup

func SortJobGroups(groups JobGroups) (JobGroups, error) {
	result := make(JobGroups, 0, len(groups))

	blocked := make(map[*JobGroup]struct{})
	mustHave := make(map[*JobGroup]struct{})

	for _, group := range groups {
		blocked[group] = struct{}{}
		if group.After == nil {
			result = append(result, group)
		} else {
			mustHave[group] = struct{}{}
			blocked[group.After] = struct{}{}
		}
	}

	for {
		move := false
		for group := range blocked {
			if _, ok := blocked[group.After]; ok {
				continue
			}

			if _, ok := mustHave[group]; ok {
				result = append(result, group)
				delete(mustHave, group)
			}

			delete(blocked, group)
			move = true
			break
		}

		if len(mustHave) == 0 {
			break
		}
		if !move {
			return nil, fmt.Errorf("The are groups in circular positionaning")
		}
	}

	return result, nil
}
