package syntax

import (
	"go/token"
)

type (
	GetFromTypeExpr struct{ typ TypeExpr }
)

func (g GetFromTypeExpr) Position() token.Position {
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
	case *UnexpectedExpr:
		return x.Position
	default:
		panic("unknown type")
	}
}

func (t *Type) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}

func (t *Map) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}

func (t *Chan) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}

func (t *FuncType) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}

func (t *StructType) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}

func (t *InterfaceType) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}

func (t *UnexpectedExpr) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}
