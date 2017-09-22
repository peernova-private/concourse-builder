package resource

import (
	"time"

	"github.com/concourse-friends/concourse-builder/model"
	"github.com/concourse-friends/concourse-builder/project"
)

// Time resource source
type TimeSource struct {
	// Lose interval between versions
	Interval time.Duration
}

// The time resource type
var TimeResourceType = &project.ResourceType{
	// The name
	Name: "time",

	// The type
	Type: model.SystemResourceTypeName,
}

func init() {
	project.GlobalTypeRegistry.MustRegisterType(TimeResourceType)
}
