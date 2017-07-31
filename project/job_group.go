package project

import (
	"fmt"
)

type JobGroup struct {
	Name   string
	After  JobGroups
	Before JobGroups
}

var SystemGroup = &JobGroup{
	Name: "sys",
}

type JobGroups []*JobGroup

func sortJobGroupsByBefore(groups JobGroups) (JobGroups, error) {
	blocked := make(map[*JobGroup]struct{})
	mustHave := make(map[*JobGroup]struct{})

	result := make(JobGroups, 0, len(groups))
	for _, group := range groups {
		if len(group.Before) == 0 {
			result = append(JobGroups{group}, result...)
			continue
		}

		mustHave[group] = struct{}{}
		blocked[group] = struct{}{}
		for _, beforeGroup := range group.Before {
			if len(beforeGroup.Before) > 0 {
				blocked[beforeGroup] = struct{}{}
			}
		}
	}

	for {
		move := false
		for group := range blocked {
			stillBlocked := false
			for _, beforeGroup := range group.Before {
				if _, ok := blocked[beforeGroup]; ok {
					stillBlocked = true
					break
				}
			}
			if stillBlocked {
				continue
			}

			if _, ok := mustHave[group]; ok {
				result = append(JobGroups{group}, result...)
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
			return nil, fmt.Errorf("There are groups in circular positionaning")
		}
	}

	return result, nil
}

// Sorts the groups based on their after and before values.
// First it puts all groups in order based on the after argument.
// With other words every group will be behind any group that is listed
// in its after list.
// Then for the groups that positions are not based on their after argument
// will be ordered based on their before argument.
// TODO: There might be better way to do the group ordering
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

	var err error
	result, err = sortJobGroupsByBefore(result)
	if err != nil {
		return nil, err
	}

	for {
		move := false
		var batch JobGroups
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
				batch = append(batch, group)
				delete(mustHave, group)
			}

			move = true
		}

		for _, group := range batch {
			delete(blocked, group)
		}

		batch, err = sortJobGroupsByBefore(batch)
		if err != nil {
			return nil, err
		}
		result = append(result, batch...)

		if len(mustHave) == 0 {
			break
		}
		if !move {
			return nil, fmt.Errorf("There are groups in circular positionaning")
		}
	}

	return result, nil
}
