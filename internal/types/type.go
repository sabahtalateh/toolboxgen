package types

import (
	"go/token"
)

type (
	Type interface {
		typ()

		Equal(Type) bool
		Get() GetFromType
	}

	Builtin struct {
		TypeName string
		Code     string
	}

	Struct struct {
		Package      string
		TypeName     string
		TypeParams   TypeParams
		Fields       Fields
		Position     token.Position
		TypePosition token.Position
		Code         string
	}

	Interface struct {
		Package      string
		TypeName     string
		TypeParams   TypeParams
		Fields       Fields
		Position     token.Position
		TypePosition token.Position
		Code         string
	}

	TypeDef struct {
		Package      string
		TypeName     string
		TypeParams   TypeParams
		Type         TypeRef
		Position     token.Position
		TypePosition token.Position
		Code         string
	}

	TypeAlias struct {
		Package      string
		TypeName     string
		Type         TypeRef
		Position     token.Position
		TypePosition token.Position
		Code         string
	}

	TypeParam struct {
		Order    int
		Name     string
		Position token.Position
		Code     string
	}

	TypeParams []*TypeParam
)

func (t *Builtin) typ()   {}
func (t *Struct) typ()    {}
func (t *Interface) typ() {}
func (t *TypeDef) typ()   {}
func (t *TypeAlias) typ() {}
