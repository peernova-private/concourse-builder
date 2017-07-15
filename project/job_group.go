package project

import (
	"fmt"
)

type JobGroup struct {
	Name  string
	After JobGroups
}

type JobGroups []*JobGroup

func SortJobGroups(groups JobGroups) (JobGroups, error) {
	blocked := make(map[*JobGroup]struct{})
	mustHave := make(map[*JobGroup]struct{})

	result := make(JobGroups, 0, len(groups))
	for _, group := range groups {
		if len(group.After) == 0 {
			result = append(result, group)
			continue
		}

		mustHave[group] = struct{}{}
		blocked[group] = struct{}{}
		for _, afterGroup := range group.After {
			if len(afterGroup.After) > 0 {
				blocked[afterGroup] = struct{}{}
			}
		}
	}

	for {
		move := false
		for group := range blocked {
			stillBlocked := false
			for _, afterGroup := range group.After {
				if _, ok := blocked[afterGroup]; ok {
					stillBlocked = true
					break
				}
			}
			if stillBlocked {
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
