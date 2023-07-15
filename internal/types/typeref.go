package types

import (
	"go/token"
)

type (
	TypeRef interface {
		typRef()
	}

	BuiltinRef struct {
		Declared  string
		Modifiers Modifiers
		TypeName  string
		Position  token.Position
	}

	StructRef struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams []TypeRef
		Position   token.Position
	}

	InterfaceRef struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams []TypeRef
		Position   token.Position
	}

	TypeDefRef struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		Type       TypeRef
		TypeParams []TypeRef
		Position   token.Position
	}

	TypeAliasRef struct {
		Declared  string
		Modifiers Modifiers
		Package   string
		TypeName  string
		Type      TypeRef
		Position  token.Position
	}

	MapRef struct {
		Declared  string
		Modifiers Modifiers
		Key       TypeRef
		Value     TypeRef
		Position  token.Position
	}

	ChanRef struct {
		Declared  string
		Modifiers Modifiers
		Value     TypeRef
		Position  token.Position
	}

	FuncTypeRef struct {
		Position token.Position
	}

	TypeParamRef struct {
		Declared  string
		Original  string
		Name      string
		Modifiers Modifiers
		Position  token.Position
	}
)

func (x *BuiltinRef) typRef()   {}
func (x *StructRef) typRef()    {}
func (x *InterfaceRef) typRef() {}
func (x *TypeDefRef) typRef()   {}
func (x *MapRef) typRef()       {}
func (x *ChanRef) typRef()      {}
func (x *FuncTypeRef) typRef()  {}
func (x *TypeParamRef) typRef() {}
func (x *TypeAliasRef) typRef() {}

func TypeRefPosition(t TypeRef) token.Position {
	switch tt := t.(type) {
	case *BuiltinRef:
		return tt.Position
	case *StructRef:
		return tt.Position
	case *InterfaceRef:
		return tt.Position
	case *MapRef:
		return tt.Position
	case *ChanRef:
		return tt.Position
	case *FuncTypeRef:
		return tt.Position
	case *TypeParamRef:
		return tt.Position
	case *TypeAliasRef:
		return tt.Position
	default:
		panic("unknown type")
	}
}
