package discovery

import "github.com/sabahtalateh/toolboxgen/internal/tool/di"

type Components []*di.Component

type Tools struct {
	Components Components
}

func newTools() *Tools {
	return &Tools{Components: []*di.Component{}}
}
