package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceRegistry(t *testing.T) {

	resourceFoo := Resource{
		Name: "foo",
		Type: ResourceTypeName("git"),
	}

	resourceBar := Resource{
		Name: "Bar",
		Type: ResourceTypeName("git"),
	}

	assert.Equal(t, resourceFoo.MustHash(), resourceBar.MustHash())
}
