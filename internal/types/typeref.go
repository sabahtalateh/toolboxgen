package types

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/clone"
)

type (
	TypeRef interface {
		typRef()

		Get() GetOnRef
		Set() SetOnRef
		Equal(TypeRef) bool

		clone.Clone[TypeRef]
	}

	BuiltinRef struct {
		TypeName   string
		Modifiers  Modifiers
		Position   token.Position
		Definition *Builtin
		Declared   string
	}

	StructRef struct {
		Package    string
		TypeName   string
		Fields     Fields
		Modifiers  Modifiers
		TypeParams TypeRefs
		Position   token.Position
		Definition *Struct
		Declared   string
	}

	InterfaceRef struct {
		Package    string
		TypeName   string
		Fields     Fields
		Modifiers  Modifiers
		TypeParams TypeRefs
		Position   token.Position
		Interface  *Interface
		Declared   string
	}

	TypeDefRef struct {
		Package    string
		TypeName   string
		Type       TypeRef
		Modifiers  Modifiers
		TypeParams TypeRefs
		Position   token.Position
		Definition *TypeDef
		Declared   string
	}

	TypeAliasRef struct {
		Package    string
		TypeName   string
		Type       TypeRef
		Modifiers  Modifiers
		Position   token.Position
		Definition *TypeAlias
		Declared   string
	}

	MapRef struct {
		Key       TypeRef
		Value     TypeRef
		Modifiers Modifiers
		Position  token.Position
		Declared  string
	}

	ChanRef struct {
		Value     TypeRef
		Modifiers Modifiers
		Position  token.Position
		Declared  string
	}

	FuncTypeRef struct {
		Params    Fields
		Results   Fields
		Modifiers Modifiers
		Position  token.Position
		Declared  string
	}

	StructTypeRef struct {
		Fields    Fields
		Modifiers Modifiers
		Position  token.Position
		Declared  string
	}

	InterfaceTypeRef struct {
		Fields    Fields
		Modifiers Modifiers
		Position  token.Position
		Declared  string
	}

	TypeParamRef struct {
		Order      int
		Modifiers  Modifiers
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
