package types

import (
	"go/token"
)

type (
	GetFromType struct{ typ Type }
)

func (g GetFromType) Position() token.Position {
	switch tt := g.typ.(type) {
	case *Builtin:
		return token.Position{}
	case *Struct:
		return tt.Position
	case *Interface:
		return tt.Position
	case *TypeDef:
		return tt.Position
	case *TypeAlias:
		return tt.Position
	default:
		panic("unknown type")
	}
}

func (g GetFromType) TypePosition() token.Position {
	switch tt := g.typ.(type) {
	case *Builtin:
		return token.Position{}
	case *Struct:
		return tt.TypePosition
	case *Interface:
		return tt.TypePosition
	case *TypeDef:
		return tt.TypePosition
	case *TypeAlias:
		return tt.TypePosition
	default:
		panic("unknown type")
	}
}

func (t *Builtin) Get() GetFromType {
	return GetFromType{typ: t}
}

func (t *Struct) Get() GetFromType {
	return GetFromType{typ: t}
}

func (t *Interface) Get() GetFromType {
	return GetFromType{typ: t}
}

func (t *TypeDef) Get() GetFromType {
	return GetFromType{typ: t}
}

func (t *TypeAlias) Get() GetFromType {
	return GetFromType{typ: t}
}

type (
	GetFromExpr struct{ expr TypeExpr }
)

func (g GetFromExpr) Position() token.Position {
	switch r := g.expr.(type) {
	case *BuiltinExpr:
		return r.Position
	case *StructExpr:
		return r.Position
	case *InterfaceExpr:
		return r.Position
	case *TypeDefExpr:
		return r.Position
	case *TypeAliasExpr:
		return r.Position
	case *MapExpr:
		return r.Position
	case *ChanExpr:
		return r.Position
	case *FuncTypeExpr:
		return r.Position
	case *StructTypeExpr:
		return r.Position
	case *TypeArgExpr:
		return r.Position
	default:
		panic("unknown type")
	}
}

func (g GetFromExpr) Modifiers() Modifiers {
	switch r := g.expr.(type) {
	case *BuiltinExpr:
		return r.Modifiers
	case *StructExpr:
		return r.Modifiers
	case *InterfaceExpr:
		return r.Modifiers
	case *TypeDefExpr:
		return r.Modifiers
	case *TypeAliasExpr:
		return r.Modifiers
	case *MapExpr:
		return r.Modifiers
	case *ChanExpr:
		return r.Modifiers
	case *FuncTypeExpr:
		return r.Modifiers
	case *StructTypeExpr:
		return r.Modifiers
	case *InterfaceTypeExpr:
		return r.Modifiers
	case *TypeArgExpr:
		return r.Modifiers
	default:
		panic("unknown type")
	}
}

func (t *BuiltinExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *StructExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *InterfaceExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *TypeDefExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *TypeAliasExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *MapExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *ChanExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *FuncTypeExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *StructTypeExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *InterfaceTypeExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

func (t *TypeArgExpr) Get() GetFromExpr {
	return GetFromExpr{expr: t}
}

type (
	SetOnExpr struct{ expr TypeExpr }
)

func (s SetOnExpr) Modifiers(m Modifiers) {
	switch r := s.expr.(type) {
	case *BuiltinExpr:
		r.Modifiers = m
	case *StructExpr:
		r.Modifiers = m
	case *InterfaceExpr:
		r.Modifiers = m
	case *TypeDefExpr:
		r.Modifiers = m
	case *TypeAliasExpr:
		r.Modifiers = m
	case *MapExpr:
		r.Modifiers = m
	case *ChanExpr:
		r.Modifiers = m
	case *FuncTypeExpr:
		r.Modifiers = m
	case *StructTypeExpr:
		r.Modifiers = m
	case *InterfaceTypeExpr:
		r.Modifiers = m
	case *TypeArgExpr:
		r.Modifiers = m
	default:
		panic("unknown type")
	}
}

func (t *BuiltinExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *StructExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *InterfaceExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *TypeDefExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *TypeAliasExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *MapExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *ChanExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *FuncTypeExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *StructTypeExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *InterfaceTypeExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}

func (t *TypeArgExpr) Set() SetOnExpr {
	return SetOnExpr{expr: t}
}
