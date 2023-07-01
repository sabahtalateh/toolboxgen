package tool

import (
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"go/token"
)

type TypeParamRef struct {
	Name     string
	Mods     []syntax.TypeRefModifier
	Position token.Position
}

func (s *TypeParamRef) typRef() {
}

func (s *TypeParamRef) Equals(t TypeRef) bool {
	switch t2 := t.(type) {
	case *TypeParamRef:
		if !syntax.ModifiersEquals(s.Modifiers(), t2.Modifiers()) {
			return false
		}

		if s.Name != t2.Name {
			return false
		}

		return true
	default:
		return false
	}
}

func (s *TypeParamRef) Modifiers() []syntax.TypeRefModifier {
	return s.Mods
}
