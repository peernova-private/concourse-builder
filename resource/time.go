package resource

import (
	"time"

	"github.com/concourse-friends/concourse-builder/model"
)

// Time resource source
type TimeSource struct {
	// Lose interval between versions
	Interval time.Duration
}

// The time resource type
var TimeResourceType = &model.ResourceType{
	// The name
	Name: "time",

	// The type
	Type: SystemResourceTypeName,
}

func init() {
	GlobalTypeRegistry.MustRegisterType(TimeResourceType)
}
