package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortJobGroups(t *testing.T) {
	a := &JobGroup{
		Name: "A",
	}
	b := &JobGroup{
		Name:  "B",
		After: a,
	}

	ab := JobGroups{a, b}
	sorted, err := SortJobGroups(ab)
	assert.NoError(t, err)
	assert.EqualValues(t, ab, sorted)

	ba := JobGroups{b, a}
	sorted, err = SortJobGroups(ba)
	assert.NoError(t, err)
	assert.EqualValues(t, ab, sorted)

	// create circle
	a.After = b
	sorted, err = SortJobGroups(ba)
	assert.Error(t, err)
}
