package inspect

import (
	"fmt"
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func TypeRef(ctx Context, t types.TypeRef) string {
	switch tt := t.(type) {
	case *types.BuiltinRef:
		return BuiltinRef(tt)
	case *types.StructRef:
		return StructRef(ctx, tt)
	case *types.InterfaceRef:
		return InterfaceRef(ctx, tt)
	case *types.TypeDefRef:
		return TypeDefRef(ctx, tt)
	case *types.TypeAliasRef:
		return TypeAliasRef(ctx, tt)
	case *types.MapRef:
		return MapRef(ctx, tt)
	case *types.ChanRef:
		return ChanRef(ctx, tt)
	case *types.FuncTypeRef:
		return FuncTypeRef(ctx, tt)
	case *types.StructTypeRef:
		return StructTypeRef(ctx, tt)
	case *types.InterfaceTypeRef:
		return InterfaceTypeRef(ctx, tt)
	case *types.TypeParamRef:
		return TypeParamRef(tt)
	default:
		panic("unknown type reference")
	}
}

func BuiltinRef(t *types.BuiltinRef) string {
	return Modifiers(t.Modifiers) + t.TypeName
}

func StructRef(ctx Context, t *types.StructRef) string {
	out := Modifiers(t.Modifiers) + typeID(ctx, t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := TypeRefs(ctx, t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func InterfaceRef(ctx Context, t *types.InterfaceRef) string {
	out := Modifiers(t.Modifiers) + typeID(ctx, t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := TypeRefs(ctx, t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func TypeDefRef(ctx Context, t *types.TypeDefRef) string {
	out := Modifiers(t.Modifiers) + typeID(ctx, t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := TypeRefs(ctx, t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func TypeAliasRef(ctx Context, t *types.TypeAliasRef) string {
	return Modifiers(t.Modifiers) + typeID(ctx, t.Package, t.TypeName)
}

func MapRef(ctx Context, t *types.MapRef) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("map[%s]", TypeRef(ctx, t.Key))
	out += fmt.Sprintf("%s", TypeRef(ctx, t.Value))

	return out
}

func ChanRef(ctx Context, t *types.ChanRef) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("chan %s", TypeRef(ctx, t.Value))

	return out
}

func FuncTypeRef(ctx Context, t *types.FuncTypeRef) string {
	out := Modifiers(t.Modifiers) + "func ("
	out += strings.Join(Fields(ctx, t.Params), ", ") + ")"

	if len(t.Results) > 0 {
		results := Fields(ctx, t.Results)
		out += " "
		if len(results) > 1 {
			out += "("
		}
		out += strings.Join(results, ", ")
		if len(results) > 1 {
			out += ")"
		}
	}

	return out
}

func StructTypeRef(ctx Context, t *types.StructTypeRef) string {
	out := Modifiers(t.Modifiers) + "struct{"

	if len(t.Fields) > 0 {
		fields := Fields(ctx, t.Fields)
		for _, field := range fields {
			out += " " + field + ";"
		}
		out = strings.TrimSuffix(out, ";")
		out += " "
	}

	out += "}"

	return out
}

func InterfaceTypeRef(ctx Context, t *types.InterfaceTypeRef) string {
	out := Modifiers(t.Modifiers) + "interface{"

	if len(t.Fields) > 0 {
		fields := Fields(ctx, t.Fields)
		for _, field := range fields {
			out += " " + field + ";"
		}
		out = strings.TrimSuffix(out, ";")
		out += " "
	}

	out += "}"

	return out
}

func TypeParamRef(t *types.TypeParamRef) string {
	return fmt.Sprintf("%s%s", Modifiers(t.Modifiers), t.Name)
}

func Field(ctx Context, f *types.Field) string {
	name := ""
	if f.Name != "" {
		name = f.Name + " "
	}
	typeOut := TypeRef(ctx, f.Type)
	return name + typeOut
}

func Fields(ctx Context, f types.Fields) []string {
	return slices.Map(f, func(el *types.Field) string { return Field(ctx, el) })
}

func TypeRefs(ctx Context, t types.TypeRefs) []string {
	return slices.Map(t, func(el types.TypeRef) string { return TypeRef(ctx, el) })
}
