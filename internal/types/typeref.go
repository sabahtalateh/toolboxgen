package types

import (
	"go/token"
)

type (
	TypeRef interface {
		typRef()
		Equal(TypeRef) bool
	}

	BuiltinRef struct {
		Declared   string
		Modifiers  Modifiers
		TypeName   string
		Definition *Builtin
		Position   token.Position
	}

	StructRef struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams []TypeRef
		Fields     Fields
		Definition *Struct
		Position   token.Position
	}

	InterfaceRef struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams []TypeRef
		Definition *Interface
		Position   token.Position
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
		Declared  string
		Modifiers Modifiers
		Params    Fields
		Results   Fields
		Position  token.Position
	}

	StructTypeRef struct {
		Declared  string
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
	}

	TypeParamRef struct {
		Declared  string
		Original  string
		Name      string
		Modifiers Modifiers
		Position  token.Position
	}

	TypeDefRef struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		Type       TypeRef
		TypeParams []TypeRef
		Definition *TypeDef
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
)

func (t *BuiltinRef) typRef()    {}
func (t *StructRef) typRef()     {}
func (t *InterfaceRef) typRef()  {}
func (t *MapRef) typRef()        {}
func (t *ChanRef) typRef()       {}
func (t *FuncTypeRef) typRef()   {}
func (t *StructTypeRef) typRef() {}
func (t *TypeParamRef) typRef()  {}
func (t *TypeDefRef) typRef()    {}
func (t *TypeAliasRef) typRef()  {}

func (t *BuiltinRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *BuiltinRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *StructRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *StructRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *InterfaceRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *InterfaceRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *MapRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *MapRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Key.Equal(tt2.Key) {
			return false
		}

		if !t.Value.Equal(tt2.Value) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *ChanRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *ChanRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Value.Equal(tt2.Value) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *FuncTypeRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *FuncTypeRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Params.Equal(tt2.Params) {
			return false
		}

		if !t.Results.Equal(tt2.Results) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *StructTypeRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *StructTypeRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Fields.Equal(tt2.Fields) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeParamRef) Equal(t2 TypeRef) bool {
	return false
}

func (t *TypeDefRef) Equal(t2 TypeRef) bool {
	return false
}

func (t *TypeAliasRef) Equal(t2 TypeRef) bool {
	return false
}
