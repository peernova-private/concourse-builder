package test

import (
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ContextDiff(expected, actual string) string {
	diff := difflib.ContextDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		FromFile: "expected",
		ToFile:   "actual",
		Context:  20,
		Eol:      "\n",
	}
	result, _ := difflib.GetContextDiffString(diff)
	return result
}

func AssertEqual(t *testing.T, expected string, actual string) {
	diff := ContextDiff(expected, actual)
	assert.True(t, diff == "", diff)
}

func RequireEqual(t *testing.T, expected string, actual string) {
	diff := ContextDiff(expected, actual)
	require.True(t, diff == "", diff)
}
