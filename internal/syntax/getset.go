package syntax

import (
	"go/token"
)

type (
	GetFromTypeRef struct{ typ TypeRef }
)

func (g GetFromTypeRef) Position() token.Position {
	switch x := g.typ.(type) {
	case *Type:
		return x.Position
	case *Map:
		return x.Position
	case *Chan:
		return x.Position
	case *FuncType:
		return x.Position
	case *StructType:
		return x.Position
	case *UnknownExpr:
		return x.Position
	default:
		panic("unknown type")
	}
}

func (t *Type) Get() GetFromTypeRef {
	return GetFromTypeRef{typ: t}
}

func (t *Map) Get() GetFromTypeRef {
	return GetFromTypeRef{typ: t}
}

func (t *Chan) Get() GetFromTypeRef {
	return GetFromTypeRef{typ: t}
}

func (t *FuncType) Get() GetFromTypeRef {
	return GetFromTypeRef{typ: t}
}

func (t *StructType) Get() GetFromTypeRef {
	return GetFromTypeRef{typ: t}
}

func (t *InterfaceType) Get() GetFromTypeRef {
	return GetFromTypeRef{typ: t}
}

func (t *UnknownExpr) Get() GetFromTypeRef {
	return GetFromTypeRef{typ: t}
}
