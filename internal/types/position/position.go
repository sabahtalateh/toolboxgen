package position

import (
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func OfType(t types.Type) token.Position {
	switch tt := t.(type) {
	case *types.Builtin:
		return token.Position{}
	case *types.Struct:
		return tt.Position
	case *types.Interface:
		return tt.Position
	case *types.TypeDef:
		return tt.Position
	case *types.TypeAlias:
		return tt.Position
	default:
		panic("unknown type")
	}
}

func OfTypeRef(t types.TypeRef) token.Position {
	switch tt := t.(type) {
	case *types.BuiltinRef:
		return tt.Position
	case *types.StructRef:
		return tt.Position
	case *types.InterfaceRef:
		return tt.Position
	case *types.TypeDefRef:
		return tt.Position
	case *types.TypeAliasRef:
		return tt.Position
	case *types.MapRef:
		return tt.Position
	case *types.ChanRef:
		return tt.Position
	case *types.FuncTypeRef:
		return tt.Position
	case *types.StructTypeRef:
		return tt.Position
	case *types.TypeParamRef:
		return tt.Position
	default:
		panic("unknown type")
	}
}
