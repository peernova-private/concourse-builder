package resource

import "github.com/concourse-friends/concourse-builder/model"

// Concourse supports several types without providing the type.
// The system resource type name indicates that.
var SystemResourceTypeName = model.ResourceTypeTypeName("system")
