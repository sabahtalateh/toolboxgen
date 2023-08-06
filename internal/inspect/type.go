package inspect

import (
	"fmt"
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (i *Inspect) Type(t types.Type) map[string]any {
	switch tt := t.(type) {
	case *types.Builtin:
		return i.Builtin(tt)
	case *types.Struct:
		return i.Struct(tt)
	case *types.Interface:
		return i.Interface(tt)
	case *types.TypeDef:
		return i.TypeDef(tt)
	case *types.TypeAlias:
		return i.TypeAlias(tt)
	default:
		panic("unknown type")
	}
}

func (i *Inspect) Builtin(t *types.Builtin) map[string]any {
	return map[string]any{"type " + t.TypeName: map[string]any{"builtin": t.TypeName}}
}

func (i *Inspect) Struct(t *types.Struct) map[string]any {
	res := map[string]any{"struct": i.composeType(t.Package, t.TypeName, t.TypeParams)}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return map[string]any{"type " + t.TypeName: res}
}

func (i *Inspect) Interface(t *types.Interface) map[string]any {
	res := map[string]any{"interface": i.composeType(t.Package, t.TypeName, t.TypeParams)}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return map[string]any{"type " + t.TypeName: res}
}

func (i *Inspect) TypeDef(t *types.TypeDef) map[string]any {
	res := map[string]any{
		"typedef": i.composeType(t.Package, t.TypeName, t.TypeParams) + " " + i.typeExpr(t.Type),
	}

	res["intro"] = i.TypeExpr(t.Type)

	return map[string]any{"type " + t.TypeName: res}
}

func (i *Inspect) TypeAlias(t *types.TypeAlias) map[string]any {
	res := map[string]any{"typealias": i.composeType(t.Package, t.TypeName, nil) + " = " + i.typeExpr(t.Type)}

	res["intro"] = i.TypeExpr(t.Type)

	return map[string]any{"type " + t.TypeName: res}
}

func (i *Inspect) TypeParam(t *types.TypeParam) string {
	return fmt.Sprintf("%s", t.Name)
}

func (i *Inspect) TypeParams(t types.TypeParams) []string {
	return slices.Map(t, func(el *types.TypeParam) string { return i.TypeParam(el) })
}

func (i *Inspect) composeType(Package, typeName string, params types.TypeParams) string {
	Type := i.typeID(Package, typeName)
	if len(params) != 0 {
		Type += fmt.Sprintf("[%s]", strings.Join(i.TypeParams(params), ", "))
	}
	return Type
}
