package component

import "github.com/sabahtalateh/toolboxgen/internal/discovery/tool"

type Constructor struct {
	Type      tool.Type
	Name      string
	WithError bool
}

func (r *Constructor) Tool() {
}

type Provider struct {
	Type      tool.Type
	Name      string
	WithError bool
}

func (r *Provider) Tool() {
}
