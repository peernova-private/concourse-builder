package library

type Concourse struct {
	// The URL of the concourse instance
	URL string

	// True if use SSL with insecure certificate
	Insecure bool

	// The team to use
	Team string

	// User to authenticate
	User string

	// Password for the user
	Password string

	// Authentication token - alternative to user/password authentication
	Token string
}

func (c *Concourse) PublicAccessEnvironment(environment map[string]interface{}) {
	environment["CONCOURSE_URL"] = c.URL
	if c.Insecure {
		environment["INSECURE"] = "true"
	}
	if c.Team != "" {
		environment["CONCOURSE_TEAM"] = c.Team
	}
}

func (c *Concourse) Environment(environment map[string]interface{}) {
	c.PublicAccessEnvironment(environment)

	if c.User != "" && c.Password != "" {
		environment["CONCOURSE_USER"] = c.User
		environment["CONCOURSE_PASSWORD"] = c.Password
	}

	if c.Token != "" {
		environment["CONCOURSE_TOKEN"] = c.Token
	}
}
