package types

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/clone"
)

type (
	TypeRef interface {
		typRef()

		clone.Clone[TypeRef]
		Get() GetOnRef
		Set() SetOnRef
		Equal(TypeRef) bool
	}

	BuiltinRef struct {
		Modifiers  Modifiers
		TypeName   string
		Position   token.Position
		Definition *Builtin
		Declared   string
	}

	StructRef struct {
		Modifiers  Modifiers
		TypeParams TypeRefs
		Package    string
		TypeName   string
		Fields     Fields
		Position   token.Position
		Definition *Struct
		Declared   string
	}

	InterfaceRef struct {
		Modifiers  Modifiers
		TypeParams TypeRefs
		Package    string
		TypeName   string
		Methods    Fields
		Position   token.Position
		Definition *Interface
		Declared   string
	}

	TypeDefRef struct {
		Modifiers  Modifiers
		TypeParams TypeRefs
		Package    string
		TypeName   string
		Type       TypeRef
		Position   token.Position
		Definition *TypeDef
		Declared   string
	}

	TypeAliasRef struct {
		Modifiers  Modifiers
		Package    string
		TypeName   string
		Type       TypeRef
		Position   token.Position
		Definition *TypeAlias
		Declared   string
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
		Modifiers  Modifiers
		Order      int
		Name       string
		Position   token.Position
		Definition *TypeParam
		Declared   string
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
