package types

import (
	"go/token"
)

type (
	Type interface {
		typ()
	}

	Builtin struct {
		Declared string
		TypeName string
	}

	Struct struct {
		Declared   string
		Package    string
		TypeName   string
		TypeParams TypeParams
		Position   token.Position
	}

	Interface struct {
		Declared   string
		Package    string
		TypeName   string
		TypeParams TypeParams
		Position   token.Position
	}

	TypeDef struct {
		Declared   string
		Package    string
		TypeName   string
		TypeParams TypeParams
		Type       TypeRef
		Position   token.Position
	}

	TypeAlias struct {
		Declared string
		Package  string
		TypeName string
		Type     TypeRef
		Position token.Position
	}
)

func (x *Builtin) typ()   {}
func (x *Struct) typ()    {}
func (x *Interface) typ() {}
func (x *TypeDef) typ()   {}
func (x *TypeAlias) typ() {}

func TypePosition(t Type) token.Position {
	switch tt := t.(type) {
	case *Builtin:
		return token.Position{}
	case *Struct:
		return tt.Position
	case *Interface:
		return tt.Position
	case *TypeDef:
		return tt.Position
	case *TypeAlias:
		return tt.Position
	default:
		panic("unknown type")
	}
}
