package inspect

import (
	"fmt"
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func Type(ctx Context, t types.Type) string {
	switch tt := t.(type) {
	case *types.Builtin:
		return Builtin(tt)
	case *types.Struct:
		return Struct(ctx, tt)
	case *types.Interface:
		return Interface(ctx, tt)
	case *types.TypeDef:
		return TypeDef(ctx, tt)
	case *types.TypeAlias:
		return TypeAlias(ctx, tt)
	default:
		panic("unknown type")
	}
}

func Builtin(t *types.Builtin) string {
	return t.TypeName
}

func Struct(ctx Context, t *types.Struct) string {
	out := typeBlock(ctx, t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		typeParamsOut := TypeParams(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	out += " struct {"
	fields := Fields(ctx, t.Fields)
	if len(fields) != 0 {
		for _, field := range fields {
			out += "\n\t" + field
		}
		out += "\n"
	}

	out += "}"

	return out
}

func Interface(ctx Context, t *types.Interface) string {
	out := typeBlock(ctx, t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		typeParamsOut := TypeParams(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	out += " interface {"
	methods := Fields(ctx, t.Fields)

	if len(methods) > 0 {
		for _, method := range methods {
			out += "\n\t" + method
		}
		out += "\n"
	}

	out += "}"

	return out
}

func TypeDef(ctx Context, t *types.TypeDef) string {
	out := typeBlock(ctx, t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		typeParamsOut := TypeParams(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	out += " " + TypeRef(ctx, t.Type)

	return out
}

func TypeAlias(ctx Context, t *types.TypeAlias) string {
	out := typeBlock(ctx, t.Package, t.TypeName)
	out += " = "
	out += TypeRef(ctx, t.Type)

	return out
}

func TypeParam(t *types.TypeParam) string {
	return fmt.Sprintf("%s", t.Name)
}

func TypeParams(t types.TypeParams) []string {
	return slices.Map(t, func(el *types.TypeParam) string { return TypeParam(el) })
}
