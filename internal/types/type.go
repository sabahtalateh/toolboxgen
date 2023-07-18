package types

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/maps"
)

type (
	Type interface {
		typ()
		Equal(Type) bool
	}

	Builtin struct {
		Declared string
		TypeName string
	}

	Struct struct {
		Declared   string
		Package    string
		TypeName   string
		TypeParams TypeParams
		Fields     Fields
		Position   token.Position
	}

	Interface struct {
		Declared   string
		Package    string
		TypeName   string
		TypeParams TypeParams
		Methods    Fields
		Position   token.Position
	}

	TypeDef struct {
		Declared   string
		Package    string
		TypeName   string
		TypeParams TypeParams
		Type       TypeRef
		Position   token.Position
	}

	TypeAlias struct {
		Declared string
		Package  string
		TypeName string
		Type     TypeRef
		Position token.Position
	}

	TypeParam struct {
		Declared string
		Original string
		Name     string
		Order    int
		Position token.Position
	}

	TypeParams []*TypeParam
)

func (t *Builtin) typ()   {}
func (t *Struct) typ()    {}
func (t *Interface) typ() {}
func (t *TypeDef) typ()   {}
func (t *TypeAlias) typ() {}

func (t *Builtin) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *Builtin:
		return t.TypeName == tt2.TypeName
	default:
		return false
	}
}

func (t *Struct) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *Struct:
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

func (t *Interface) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *Interface:
		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		methods1 := maps.FromSlice(t.Methods, func(s *Field) (string, *Field) { return s.Name, s })
		methods2 := maps.FromSlice(tt2.Methods, func(s *Field) (string, *Field) { return s.Name, s })

		if len(methods1) != len(methods2) {
			return false
		}

		for k, v := range methods1 {
			v2, ok := methods2[k]
			if !ok {
				return false
			}
			if !v.Equal(v2) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func (t *TypeDef) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *TypeDef:
		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		if !t.Type.Equal(tt2.Type) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeAlias) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *TypeAlias:
		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.Type.Equal(tt2.Type) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeParam) Equal(t2 *TypeParam) bool {
	return t.Name == t2.Name && t.Order == t2.Order
}

func (t TypeParams) Equal(t2 TypeParams) bool {
	if len(t) != len(t2) {
		return false
	}

	for i, param := range t {
		if !param.Equal(t2[i]) {
			return false
		}
	}

	return true
}
