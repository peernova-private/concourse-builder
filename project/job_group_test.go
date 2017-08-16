package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortJobGroupsEmpty(t *testing.T) {
	empty := JobGroups{}
	sorted, err := SortJobGroups(empty)
	assert.NoError(t, err)
	assert.EqualValues(t, empty, sorted)
}

func TestSortJobGroupsA(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}

	slice := JobGroups{a}
	sorted, err := SortJobGroups(slice)
	assert.NoError(t, err)
	assert.EqualValues(t, slice, sorted)
}

func TestSortJobGroupsACircleAfter(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}
	a.After = JobGroups{a}

	slice := JobGroups{a}
	_, err := SortJobGroups(slice)
	assert.Error(t, err)
}

func TestSortJobGroupsACircleBefore(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}
	a.Before = JobGroups{a}

	slice := JobGroups{a}
	_, err := SortJobGroups(slice)
	assert.Error(t, err)
}

func TestSortJobGroupsAB(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}
	b := &JobGroup{
		Name:  "B",
		After: JobGroups{a},
	}

	ab := JobGroups{a, b}
	sorted, err := SortJobGroups(ab)
	assert.NoError(t, err)
	assert.EqualValues(t, ab, sorted)
}

func TestSortJobGroupsBA(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}
	b := &JobGroup{
		Name:  "B",
		After: JobGroups{a},
	}

	ab := JobGroups{a, b}
	ba := JobGroups{b, a}
	sorted, err := SortJobGroups(ba)
	assert.NoError(t, err)
	assert.EqualValues(t, ab, sorted)
}

func TestSortJobGroupsABCircle(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}
	b := &JobGroup{
		Name:  "B",
		After: JobGroups{a},
	}
	a.After = JobGroups{b}

	ab := JobGroups{a, b}
	_, err := SortJobGroups(ab)
	assert.Error(t, err)
}

func TestSortJobGroupsCBA(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}
	b := &JobGroup{
		Name:  "B",
		After: JobGroups{a},
	}
	c := &JobGroup{
		Name:  "C",
		After: JobGroups{b},
	}

	abc := JobGroups{a, b, c}
	cba := JobGroups{c, b, a}
	sorted, err := SortJobGroups(cba)
	assert.NoError(t, err)
	assert.EqualValues(t, abc, sorted)
}
