package types

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/clone"
)

type (
	TypeRef interface {
		typRef()

		Get() GetFromRef
		Set() SetOnRef
		Equal(TypeRef) bool

		clone.Clone[TypeRef]
	}

	BuiltinRef struct {
		TypeName   string
		Modifiers  Modifiers
		Position   token.Position
		Definition *Builtin
		Code       string
	}

	StructRef struct {
		Package    string
		TypeName   string
		Fields     Fields
		Modifiers  Modifiers
		TypeParams TypeRefs
		Position   token.Position
		Definition *Struct
		Code       string
	}

	InterfaceRef struct {
		Package    string
		TypeName   string
		Fields     Fields
		Modifiers  Modifiers
		TypeParams TypeRefs
		Position   token.Position
		Definition *Interface
		Code       string
	}

	TypeDefRef struct {
		Package    string
		TypeName   string
		Type       TypeRef
		Modifiers  Modifiers
		TypeParams TypeRefs
		Position   token.Position
		Definition *TypeDef
		Code       string
	}

	TypeAliasRef struct {
		Package    string
		TypeName   string
		Type       TypeRef
		Modifiers  Modifiers
		Position   token.Position
		Definition *TypeAlias
		Code       string
	}

	MapRef struct {
		Key       TypeRef
		Value     TypeRef
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	ChanRef struct {
		Value     TypeRef
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	FuncTypeRef struct {
		Params    Fields
		Results   Fields
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	StructTypeRef struct {
		Fields    Fields
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	InterfaceTypeRef struct {
		Fields    Fields
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	TypeParamRef struct {
		Order      int
		Modifiers  Modifiers
		Name       string
		Position   token.Position
		Definition *TypeParam
		Code       string
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
