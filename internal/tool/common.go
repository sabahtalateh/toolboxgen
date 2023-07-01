package tool

import (
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"go/token"
)

func Position(t TypeRef) token.Position {
	switch x := t.(type) {
	case *BuiltinRef:
		return x.Position
	case *StructRef:
		return x.Position
	case *InterfaceRef:
		return x.Position
	case *TypeParamRef:
		return x.Position
	default:
		panic("not implemented")
	}
}

func prependModifiers(mm []syntax.TypeRefModifier, t TypeRef) TypeRef {
	switch tt := t.(type) {
	case *StructRef:
		tt.Mods = append(mm, tt.Mods...)
	case *InterfaceRef:
		tt.Mods = append(mm, tt.Mods...)
	case *BuiltinRef:
		tt.Mods = append(mm, tt.Mods...)
	case *TypeParamRef:
		tt.Mods = append(mm, tt.Mods...)
	}

	return t
}
