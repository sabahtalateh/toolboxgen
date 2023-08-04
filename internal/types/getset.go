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
	GetFromRef struct{ ref TypeRef }
)

func (g GetFromRef) Position() token.Position {
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

func (g GetFromRef) Modifiers() Modifiers {
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
	case *InterfaceTypeRef:
		return r.Modifiers
	case *TypeParamRef:
		return r.Modifiers
	default:
		panic("unknown type")
	}
}

func (t *BuiltinRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *StructRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *InterfaceRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *TypeDefRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *TypeAliasRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *MapRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *ChanRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *FuncTypeRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *StructTypeRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *InterfaceTypeRef) Get() GetFromRef {
	return GetFromRef{ref: t}
}

func (t *TypeParamRef) Get() GetFromRef {
	return GetFromRef{ref: t}
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
	case *InterfaceTypeRef:
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
