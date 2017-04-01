package resource

import "github.com/ggeorgiev/concourse-builder/model"

// Concourse supports several types without providing the type.
// The system resource type name indicates that.
var SystemResourceTypeName = model.ResourceTypeTypeName("system")
