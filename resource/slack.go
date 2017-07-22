package resource

import (
	"github.com/concourse-friends/concourse-builder/model"
)

// The slack resource source
type SlackSource struct {
	// Slack access URL
	URL string
}

// Slack resource type
var SlackResourceType = &model.ResourceType{
	// The name
	Name: "slack-notification",

	// The type
	Type: ImageResourceType.Type,

	// The official resource deployment
	Source: &ImageSource{
		// The image repository
		Repository: "cfcommunity/slack-notification-resource",
	},
}

func init() {
	GlobalRegistry.MustRegisterType(SlackResourceType)
}
