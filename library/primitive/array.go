package primitive

import (
	"fmt"
	"strings"

	"github.com/concourse-friends/concourse-builder/project"
	"github.com/davecgh/go-spew/spew"
)

type Array []interface{}

func (a Array) Value() string {
	values := []string{}
	for _, i := range a {
		if item, ok := i.(string); ok {
			values = append(values, item)
		} else if item, ok := i.(project.IEnvironmentValue); ok {
			values = append(values, item.Value())
		} else if item, ok := i.(fmt.Stringer); ok {
			values = append(values, item.String())
		} else {
			panic(spew.Sdump(i))
		}
	}

	return strings.Join(values, " ")
}

func (a Array) OutputNames() []string {
	var names []string
	for _, i := range a {
		if item, ok := i.(project.ITaskInput); ok {
			names = append(names, item.OutputNames()...)
		}
	}
	return names
}
