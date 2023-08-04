package mid

import (
	"go/token"
)

type (
	GetOnTypeRef struct{ typ TypeRef }
)

func (g GetOnTypeRef) Position() token.Position {
	switch x := g.typ.(type) {
	case *UnknownExpr:
		return x.Position
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
	default:
		panic("unknown type")
	}
}

func (t *Type) Get() GetOnTypeRef {
	return GetOnTypeRef{typ: t}
}

func (t *Map) Get() GetOnTypeRef {
	return GetOnTypeRef{typ: t}
}

func (t *Chan) Get() GetOnTypeRef {
	return GetOnTypeRef{typ: t}
}

func (t *FuncType) Get() GetOnTypeRef {
	return GetOnTypeRef{typ: t}
}

func (t *StructType) Get() GetOnTypeRef {
	return GetOnTypeRef{typ: t}
}

func (t *InterfaceType) Get() GetOnTypeRef {
	return GetOnTypeRef{typ: t}
}
