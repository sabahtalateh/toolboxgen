package component

import "github.com/sabahtalateh/toolboxgen/internal/discovery/tool"

type Register struct {
	Type tool.Type
}

func (r *Register) Tool() {
}

type Provider struct {
	Type tool.Type
}

func (r *Provider) Tool() {
}
