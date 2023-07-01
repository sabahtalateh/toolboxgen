package tool

import "github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"

type Tool interface {
	Tool()
}

type TypeDef interface {
	typDef()
}

type TypeRef interface {
	typRef()
	Equals(t TypeRef) bool
	Modifiers() []syntax.TypeRefModifier
}

func TypeRefFromDef(d TypeDef) TypeRef {
	switch dd := d.(type) {
	case *BuiltinDef:
		return BuiltinRefFromDef(dd)
	case *StructDef:
		return StructRefFromDef(dd)
	case *InterfaceDef:
		return InterfaceRefFromDef(dd)
	default:
		panic("not supported")
	}
}
