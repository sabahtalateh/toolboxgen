package inspect

import (
	"fmt"
	"github.com/life4/genesis/slices"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (i *Inspect) TypeExpr(t types.TypeExpr) any {
	switch tt := t.(type) {
	case *types.BuiltinExpr:
		return Modifiers(tt.Modifiers) + tt.TypeName
	case *types.StructExpr:
		return i.StructExpr(tt)
	case *types.InterfaceExpr:
		return i.InterfaceExpr(tt)
	case *types.TypeDefExpr:
		return i.TypeDefExpr(tt)
	case *types.TypeAliasExpr:
		return i.TypeAliasExpr(tt)
	case *types.MapExpr:
		return i.MapExpr(tt)
	case *types.ChanExpr:
		return i.ChanExpr(tt)
	case *types.FuncTypeExpr:
		return i.FuncTypeExpr(tt)
	case *types.StructTypeExpr:
		return i.StructTypeExpr(tt)
	case *types.InterfaceTypeExpr:
		return i.InterfaceTypeExpr(tt)
	case *types.TypeArgExpr:
		return Modifiers(tt.Modifiers) + tt.Name
	default:
		panic("unknown type")
	}
}

func (i *Inspect) typeExpr(t types.TypeExpr) string {
	switch tt := t.(type) {
	case *types.BuiltinExpr:
		return i.builtinExpr(tt)
	case *types.StructExpr:
		return i.structExpr(tt)
	case *types.InterfaceExpr:
		return i.interfaceExpr(tt)
	case *types.TypeDefExpr:
		return i.typeDefExpr(tt)
	case *types.TypeAliasExpr:
		return i.typeAliasExpr(tt)
	case *types.MapExpr:
		return i.mapExpr(tt)
	case *types.ChanExpr:
		return i.chanExpr(tt)
	case *types.FuncTypeExpr:
		return i.funcTypeExpr(tt)
	case *types.StructTypeExpr:
		return i.structTypeExpr(tt)
	case *types.InterfaceTypeExpr:
		return i.interfaceTypeExpr(tt)
	case *types.TypeArgExpr:
		return i.typeArgExpr(tt)
	default:
		panic("unknown type expression")
	}
}

func (i *Inspect) builtinExpr(t *types.BuiltinExpr) string {
	return Modifiers(t.Modifiers) + t.TypeName
}

func (i *Inspect) StructExpr(t *types.StructExpr) map[string]any {
	d := t.Definition

	res := map[string]any{"struct": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, d.TypeParams)}

	if len(t.TypeArgs) != 0 {
		res["actual"] = slices.Map(t.TypeArgs, func(el types.TypeExpr) any { return i.TypeExpr(el) })
	}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) structExpr(t *types.StructExpr) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeArgs) > 0 {
		typeArgsOut := i.typeExprs(t.TypeArgs)
		out += fmt.Sprintf("[%s]", strings.Join(typeArgsOut, ", "))
	}
	return out
}

func (i *Inspect) InterfaceExpr(t *types.InterfaceExpr) map[string]any {
	d := t.Definition

	res := map[string]any{"interface": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, d.TypeParams)}

	if len(t.TypeArgs) != 0 {
		res["actual"] = slices.Map(t.TypeArgs, func(el types.TypeExpr) any { return i.TypeExpr(el) })
	}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) interfaceExpr(t *types.InterfaceExpr) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeArgs) > 0 {
		typeArgsOut := i.typeExprs(t.TypeArgs)
		out += fmt.Sprintf("[%s]", strings.Join(typeArgsOut, ", "))
	}

	return out
}

func (i *Inspect) TypeDefExpr(t *types.TypeDefExpr) map[string]any {
	d := t.Definition

	res := map[string]any{
		"typedef": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, d.TypeParams) + " " + i.typeExpr(d.Type),
	}

	if len(t.TypeArgs) != 0 {
		res["actual"] = slices.Map(t.TypeArgs, func(el types.TypeExpr) any { return i.TypeExpr(el) })
	}

	res["intro"] = i.TypeExpr(t.Type)

	return res
}

func (i *Inspect) typeDefExpr(t *types.TypeDefExpr) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeArgs) > 0 {
		typeArgsOut := i.typeExprs(t.TypeArgs)
		out += fmt.Sprintf("[%s]", strings.Join(typeArgsOut, ", "))
	}
	return out
}

func (i *Inspect) TypeAliasExpr(t *types.TypeAliasExpr) map[string]any {
	d := t.Definition

	res := map[string]any{"typealias": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, nil) + " = " + i.typeExpr(d.Type)}

	res["intro"] = i.TypeExpr(t.Type)

	return res
}

func (i *Inspect) typeAliasExpr(t *types.TypeAliasExpr) string {
	return Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
}

func (i *Inspect) MapExpr(t *types.MapExpr) map[string]any {
	res := map[string]any{
		"map":   Modifiers(t.Modifiers) + "map[..]..",
		"key":   i.TypeExpr(t.Key),
		"value": i.TypeExpr(t.Value),
	}

	return res
}

func (i *Inspect) mapExpr(t *types.MapExpr) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("map[%s]", i.typeExpr(t.Key))
	out += fmt.Sprintf("%s", i.typeExpr(t.Value))

	return out
}

func (i *Inspect) ChanExpr(t *types.ChanExpr) map[string]any {
	res := map[string]any{
		"chan": Modifiers(t.Modifiers) + "chan ..",
		"type": i.TypeExpr(t.Value),
	}

	return res
}

func (i *Inspect) chanExpr(t *types.ChanExpr) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("chan %s", i.typeExpr(t.Value))

	return out
}

func (i *Inspect) FuncTypeExpr(t *types.FuncTypeExpr) any {
	res := map[string]any{"func": Modifiers(t.Modifiers) + "{..}"}

	if len(t.Params) != 0 {
		res["params"] = i.Fields(t.Params)
	}

	if len(t.Results) != 0 {
		res["results"] = i.Fields(t.Results)
	}

	return res
}

func (i *Inspect) funcTypeExpr(t *types.FuncTypeExpr) string {
	out := Modifiers(t.Modifiers) + "func ("
	out += strings.Join(i.fields(t.Params), ", ") + ")"

	if len(t.Results) > 0 {
		results := i.fields(t.Results)
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

func (i *Inspect) StructTypeExpr(t *types.StructTypeExpr) any {
	res := map[string]any{"struct": Modifiers(t.Modifiers) + "{..}"}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) structTypeExpr(t *types.StructTypeExpr) string {
	out := Modifiers(t.Modifiers) + "struct{"

	if len(t.Fields) > 0 {
		fields := i.fields(t.Fields)
		for _, field := range fields {
			out += " " + field + ";"
		}
		out = strings.TrimSuffix(out, ";")
		out += " "
	}

	out += "}"

	return out
}

func (i *Inspect) InterfaceTypeExpr(t *types.InterfaceTypeExpr) any {
	res := map[string]any{"interface": Modifiers(t.Modifiers) + "{..}"}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) interfaceTypeExpr(t *types.InterfaceTypeExpr) string {
	out := Modifiers(t.Modifiers) + "interface{"

	if len(t.Fields) > 0 {
		fields := i.fields(t.Fields)
		for _, field := range fields {
			out += " " + field + ";"
		}
		out = strings.TrimSuffix(out, ";")
		out += " "
	}

	out += "}"

	return out
}

func (i *Inspect) typeArgExpr(t *types.TypeArgExpr) string {
	return fmt.Sprintf("%s%s", Modifiers(t.Modifiers), t.Name)
}

func (i *Inspect) typeExprs(t types.TypeExprs) []string {
	return slices.Map(t, func(el types.TypeExpr) string { return i.typeExpr(el) })
}
