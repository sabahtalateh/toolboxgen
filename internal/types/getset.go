package types

import (
	"go/token"
)

type (
	GetOnType struct{ typ Type }
)

func (g GetOnType) Position() token.Position {
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

func (g GetOnType) TypePosition() token.Position {
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

func (t *Builtin) Get() GetOnType {
	return GetOnType{typ: t}
}

func (t *Struct) Get() GetOnType {
	return GetOnType{typ: t}
}

func (t *Interface) Get() GetOnType {
	return GetOnType{typ: t}
}

func (t *TypeDef) Get() GetOnType {
	return GetOnType{typ: t}
}

func (t *TypeAlias) Get() GetOnType {
	return GetOnType{typ: t}
}

type (
	GetOnRef struct{ ref TypeRef }
)

func (g GetOnRef) Position() token.Position {
	switch r := g.ref.(type) {
	case *BuiltinRef:
		return r.Position
	case *StructRef:
		return r.Position
	case *InterfaceRef:
		return r.Position
	case *TypeDefRef:
		return r.Position
	case *TypeAliasRef:
		return r.Position
	case *MapRef:
		return r.Position
	case *ChanRef:
		return r.Position
	case *FuncTypeRef:
		return r.Position
	case *StructTypeRef:
		return r.Position
	case *TypeParamRef:
		return r.Position
	default:
		panic("unknown type")
	}
}

func (g GetOnRef) Modifiers() Modifiers {
	switch r := g.ref.(type) {
	case *BuiltinRef:
		return r.Modifiers
	case *StructRef:
		return r.Modifiers
	case *InterfaceRef:
		return r.Modifiers
	case *TypeDefRef:
		return r.Modifiers
	case *TypeAliasRef:
		return r.Modifiers
	case *MapRef:
		return r.Modifiers
	case *ChanRef:
		return r.Modifiers
	case *FuncTypeRef:
		return r.Modifiers
	case *StructTypeRef:
		return r.Modifiers
	case *TypeParamRef:
		return r.Modifiers
	default:
		panic("unknown type")
	}
}

func (t *BuiltinRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *StructRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *InterfaceRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *TypeDefRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *TypeAliasRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *MapRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *ChanRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *FuncTypeRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *StructTypeRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *InterfaceTypeRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

func (t *TypeParamRef) Get() GetOnRef {
	return GetOnRef{ref: t}
}

type (
	SetOnRef struct{ ref TypeRef }
)

func (s SetOnRef) Modifiers(m Modifiers) {
	switch r := s.ref.(type) {
	case *BuiltinRef:
		r.Modifiers = m
	case *StructRef:
		r.Modifiers = m
	case *InterfaceRef:
		r.Modifiers = m
	case *TypeDefRef:
		r.Modifiers = m
	case *TypeAliasRef:
		r.Modifiers = m
	case *MapRef:
		r.Modifiers = m
	case *ChanRef:
		r.Modifiers = m
	case *FuncTypeRef:
		r.Modifiers = m
	case *StructTypeRef:
		r.Modifiers = m
	case *TypeParamRef:
		r.Modifiers = m
	default:
		panic("unknown type")
	}
}

func (t *BuiltinRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *StructRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *InterfaceRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *TypeDefRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *TypeAliasRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *MapRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *ChanRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *FuncTypeRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *StructTypeRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *InterfaceTypeRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}

func (t *TypeParamRef) Set() SetOnRef {
	return SetOnRef{ref: t}
}
