package tool

import "go/token"

type Tool interface {
	Tool()
}

type TypeRef interface {
	typ()
	Equals(t TypeRef) bool
}

func Position(t TypeRef) token.Position {
	switch x := t.(type) {
	case *BuiltinRef:
		return x.Position
	case *StructRef:
		return x.Position
	case *InterfaceRef:
		return x.Position
	case *TypeParamRef:
		return x.Position
	default:
		panic("not implemented")
	}
}
