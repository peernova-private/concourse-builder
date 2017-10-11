package primitive

import "regexp"

type GitRepo struct {
	// URI to the git repo
	URI string

	// Private key the allows access to the repo
	PrivateKey string
}

var githubURL = regexp.MustCompile(`^git@github.com:.*?\/(.*?).git$`)

func (gr *GitRepo) FriendlyName() string {
	if match := githubURL.FindAllStringSubmatch(gr.URI, -1); match != nil {
		return match[0][1]
	}

	// We do not know what this is, return it as is
	return gr.URI
}
