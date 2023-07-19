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
		Modifiers  Modifiers
		TypeName   string
		Definition *Builtin
		Position   token.Position
		Declared   string
	}

	StructRef struct {
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams TypeRefs
		Fields     Fields
		Definition *Struct
		Position   token.Position
		Declared   string
	}

	InterfaceRef struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams TypeRefs
		Definition *Interface
		Position   token.Position
	}

	TypeDefRef struct {
		Modifiers  Modifiers
		Package    string
		TypeName   string
		Type       TypeRef
		TypeParams TypeRefs
		Definition *TypeDef
		Position   token.Position
		Declared   string
	}

	TypeAliasRef struct {
		Modifiers Modifiers
		Package   string
		TypeName  string
		Type      TypeRef
		Position  token.Position
		Declared  string
	}

	MapRef struct {
		Modifiers Modifiers
		Key       TypeRef
		Value     TypeRef
		Position  token.Position
		Declared  string
	}

	ChanRef struct {
		Modifiers Modifiers
		Value     TypeRef
		Position  token.Position
		Declared  string
	}

	FuncTypeRef struct {
		Modifiers Modifiers
		Params    Fields
		Results   Fields
		Position  token.Position
		Declared  string
	}

	StructTypeRef struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Declared  string
	}

	InterfaceTypeRef struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Declared  string
	}

	TypeParamRef struct {
		Original  string
		Name      string
		Modifiers Modifiers
		Position  token.Position
		Declared  string
	}

	TypeRefs []TypeRef
)

func (t *BuiltinRef) typRef()       {}
func (t *StructRef) typRef()        {}
func (t *InterfaceRef) typRef()     {}
func (t *TypeDefRef) typRef()       {}
func (t *TypeAliasRef) typRef()     {}
func (t *MapRef) typRef()           {}
func (t *ChanRef) typRef()          {}
func (t *FuncTypeRef) typRef()      {}
func (t *StructTypeRef) typRef()    {}
func (t *InterfaceTypeRef) typRef() {}
func (t *TypeParamRef) typRef()     {}

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

		if !t.TypeParams.Equal(tt2.TypeParams) {
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

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeDefRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *TypeDefRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeAliasRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *TypeAliasRef:
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

func (t *InterfaceTypeRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *InterfaceTypeRef:
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
	switch tt2 := t2.(type) {
	case *TypeParamRef:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Name != tt2.Name {
			return false
		}

		return true
	default:
		return false
	}
}

func (t TypeRefs) Equal(t2 TypeRefs) bool {
	if len(t) != len(t2) {
		return false
	}

	for i, ref := range t {
		if !ref.Equal(t2[i]) {
			return false
		}
	}

	return true
}
