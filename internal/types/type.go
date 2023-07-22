package types

import (
	"go/token"
)

type (
	Type interface {
		typ()

		Equal(Type) bool
		Get() GetOnType
	}

	Builtin struct {
		TypeName string
		Declared string
	}

	Struct struct {
		TypeParams   TypeParams
		Package      string
		TypeName     string
		Fields       Fields
		Position     token.Position
		TypePosition token.Position
		Declared     string
	}

	Interface struct {
		TypeParams   TypeParams
		Package      string
		TypeName     string
		Methods      Fields
		Position     token.Position
		TypePosition token.Position
		Declared     string
	}

	TypeDef struct {
		TypeParams   TypeParams
		Package      string
		TypeName     string
		Type         TypeRef
		Position     token.Position
		TypePosition token.Position
		Declared     string
	}

	TypeAlias struct {
		Package      string
		TypeName     string
		Type         TypeRef
		Position     token.Position
		TypePosition token.Position
		Declared     string
	}

	TypeParam struct {
		Order    int
		Name     string
		Position token.Position
		Declared string
	}

	TypeParams []*TypeParam
)

func (t *Builtin) typ()   {}
func (t *Struct) typ()    {}
func (t *Interface) typ() {}
func (t *TypeDef) typ()   {}
func (t *TypeAlias) typ() {}
