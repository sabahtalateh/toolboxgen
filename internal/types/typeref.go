package types

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/clone"
)

type (
	TypeRef interface {
		typeRef()

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

func (t *BuiltinRef) typeRef()       {}
func (t *StructRef) typeRef()        {}
func (t *InterfaceRef) typeRef()     {}
func (t *TypeDefRef) typeRef()       {}
func (t *TypeAliasRef) typeRef()     {}
func (t *MapRef) typeRef()           {}
func (t *ChanRef) typeRef()          {}
func (t *FuncTypeRef) typeRef()      {}
func (t *StructTypeRef) typeRef()    {}
func (t *InterfaceTypeRef) typeRef() {}
func (t *TypeParamRef) typeRef()     {}
