package position

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/mid"
)

func OfTypeRef(r mid.TypeRef) token.Position {
	switch x := r.(type) {
	case *mid.Type:
		return x.Position
	case *mid.Map:
		return x.Position
	case *mid.Chan:
		return x.Position
	case *mid.FuncType:
		return x.Position
	case *mid.StructType:
		return x.Position
	default:
		panic("unknown type reference")
	}
}
