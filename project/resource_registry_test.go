package project

import (
	"testing"

	"github.com/concourse-friends/concourse-builder/resource"
	"github.com/stretchr/testify/assert"
)

func TestResourceRegistry(t *testing.T) {

	resourceFoo := Resource{
		Name: "foo",
		Type: resource.GitResourceType.Name,
	}

	resourceBar := Resource{
		Name: "Bar",
		Type: resource.GitResourceType.Name,
	}

	assert.Equal(t, resourceFoo.MustHash(), resourceBar.MustHash())
}
