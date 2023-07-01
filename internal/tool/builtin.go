package tool

import (
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"go/token"
)

type BuiltinDef struct {
	TypeName  string
	Modifiers []syntax.TypeRefModifier
	Position  token.Position
}

func (b BuiltinDef) typDef() {}

type BuiltinRef struct {
	TypeName string
	Mods     []syntax.TypeRefModifier
	Position token.Position
}

func BuiltinRefFromDef(d *BuiltinDef) *BuiltinRef {
	return &BuiltinRef{
		TypeName: d.TypeName,
		Mods:     d.Modifiers,
		Position: d.Position,
	}
}

func (b *BuiltinRef) typRef() {}

func (b *BuiltinRef) Equals(t TypeRef) bool {
	switch t2 := t.(type) {
	case *BuiltinRef:
		if b.TypeName != t2.TypeName {
			return false
		}

		return true
	default:
		return false
	}
}

func (b *BuiltinRef) Modifiers() []syntax.TypeRefModifier {
	return b.Mods
}

func (b *BuiltinRef) IsError() bool {
	return b.TypeName == "error"
}
