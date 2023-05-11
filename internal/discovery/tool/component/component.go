package component

import "github.com/sabahtalateh/toolboxgen/internal/discovery/tool"

type Component struct {
	Type      tool.Type
	Name      string
	WithError bool
}

func (r *Component) Tool() {
}

type Provider struct {
	Type      tool.Type
	Name      string
	WithError bool
}

func (r *Provider) Tool() {
}
