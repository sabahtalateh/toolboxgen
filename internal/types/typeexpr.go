package types

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/clone"
)

type (
	TypeExpr interface {
		typeExpr()

		clone.Clone[TypeExpr]
		Get() GetFromExpr
		Set() SetOnExpr
		Equal(TypeExpr) bool
	}

	BuiltinExpr struct {
		TypeName   string
		Modifiers  Modifiers
		Position   token.Position
		Definition *Builtin
		Code       string
	}

	StructExpr struct {
		Package    string
		TypeName   string
		Fields     Fields
		Modifiers  Modifiers
		TypeArgs   TypeExprs
		Position   token.Position
		Definition *Struct
		Code       string
	}

	InterfaceExpr struct {
		Package    string
		TypeName   string
		Fields     Fields
		Modifiers  Modifiers
		TypeArgs   TypeExprs
		Position   token.Position
		Definition *Interface
		Code       string
	}

	TypeDefExpr struct {
		Package    string
		TypeName   string
		Type       TypeExpr
		Modifiers  Modifiers
		TypeArgs   TypeExprs
		Position   token.Position
		Definition *TypeDef
		Code       string
	}

	TypeAliasExpr struct {
		Package    string
		TypeName   string
		Type       TypeExpr
		Modifiers  Modifiers
		Position   token.Position
		Definition *TypeAlias
		Code       string
	}

	MapExpr struct {
		Key       TypeExpr
		Value     TypeExpr
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	ChanExpr struct {
		Value     TypeExpr
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	FuncTypeExpr struct {
		Params    Fields
		Results   Fields
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	StructTypeExpr struct {
		Fields    Fields
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	InterfaceTypeExpr struct {
		Fields    Fields
		Modifiers Modifiers
		Position  token.Position
		Code      string
	}

	TypeArgExpr struct {
		Order      int
		Modifiers  Modifiers
		Name       string
		Position   token.Position
		Definition *TypeParam
		Code       string
	}

	TypeExprs []TypeExpr
)

func (t *BuiltinExpr) typeExpr()       {}
func (t *StructExpr) typeExpr()        {}
func (t *InterfaceExpr) typeExpr()     {}
func (t *TypeDefExpr) typeExpr()       {}
func (t *TypeAliasExpr) typeExpr()     {}
func (t *MapExpr) typeExpr()           {}
func (t *ChanExpr) typeExpr()          {}
func (t *FuncTypeExpr) typeExpr()      {}
func (t *StructTypeExpr) typeExpr()    {}
func (t *InterfaceTypeExpr) typeExpr() {}
func (t *TypeArgExpr) typeExpr()       {}
