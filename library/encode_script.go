package library

import (
	"github.com/concourse-friends/concourse-builder/project"
	"github.com/concourse-friends/concourse-builder/library/primitive"
	"fmt"
	"encoding/base64"
)

func EncodeScript(script string) (project.IRun, []interface{}) {
	bash :=  &primitive.Location{
		Volume: &primitive.Directory{
			Root: "/bin",
		},
		RelativePath: "bash",
	}

	arguments := []interface{} {
		"-c",
		fmt.Sprintf("echo %s | base64 --decode > script.sh && chmod 755 script.sh && ./script.sh",
			base64.StdEncoding.EncodeToString([]byte(script))),
	}

	return bash, arguments
}
