package inspect

import (
	"fmt"
	"github.com/life4/genesis/slices"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (i *Inspect) TypeRef(t types.TypeRef) any {
	switch tt := t.(type) {
	case *types.BuiltinRef:
		return Modifiers(tt.Modifiers) + tt.TypeName
	case *types.StructRef:
		return i.StructRef(tt)
	case *types.InterfaceRef:
		return i.InterfaceRef(tt)
	case *types.TypeDefRef:
		return i.TypeDefRef(tt)
	case *types.TypeAliasRef:
		return i.TypeAliasRef(tt)
	case *types.MapRef:
		return i.MapRef(tt)
	case *types.ChanRef:
		return i.ChanRef(tt)
	case *types.FuncTypeRef:
		return i.FuncTypeRef(tt)
	case *types.StructTypeRef:
		return i.StructTypeRef(tt)
	case *types.InterfaceTypeRef:
		return i.InterfaceTypeRef(tt)
	case *types.TypeParamRef:
		return Modifiers(tt.Modifiers) + tt.Name
	default:
		panic("unknown type")
	}
}

func (i *Inspect) typeRef(t types.TypeRef) string {
	switch tt := t.(type) {
	case *types.BuiltinRef:
		return i.builtinRef(tt)
	case *types.StructRef:
		return i.structRef(tt)
	case *types.InterfaceRef:
		return i.interfaceRef(tt)
	case *types.TypeDefRef:
		return i.typeDefRef(tt)
	case *types.TypeAliasRef:
		return i.typeAliasRef(tt)
	case *types.MapRef:
		return i.mapRef(tt)
	case *types.ChanRef:
		return i.chanRef(tt)
	case *types.FuncTypeRef:
		return i.funcTypeRef(tt)
	case *types.StructTypeRef:
		return i.structTypeRef(tt)
	case *types.InterfaceTypeRef:
		return i.interfaceTypeRef(tt)
	case *types.TypeParamRef:
		return i.typeParamRef(tt)
	default:
		panic("unknown type reference")
	}
}

func (i *Inspect) builtinRef(t *types.BuiltinRef) string {
	return Modifiers(t.Modifiers) + t.TypeName
}

func (i *Inspect) StructRef(t *types.StructRef) map[string]any {
	d := t.Definition

	res := map[string]any{"struct": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, d.TypeParams)}

	if len(t.TypeParams) != 0 {
		res["actual"] = slices.Map(t.TypeParams, func(el types.TypeRef) any { return i.TypeRef(el) })
	}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) structRef(t *types.StructRef) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := i.typeRefs(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func (i *Inspect) InterfaceRef(t *types.InterfaceRef) map[string]any {
	d := t.Definition

	res := map[string]any{"interface": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, d.TypeParams)}

	if len(t.TypeParams) != 0 {
		res["actual"] = slices.Map(t.TypeParams, func(el types.TypeRef) any { return i.TypeRef(el) })
	}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) interfaceRef(t *types.InterfaceRef) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := i.typeRefs(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}

	return out
}

func (i *Inspect) TypeDefRef(t *types.TypeDefRef) map[string]any {
	d := t.Definition

	res := map[string]any{
		"typedef": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, d.TypeParams) + " " + i.typeRef(d.Type),
	}

	if len(t.TypeParams) != 0 {
		res["actual"] = slices.Map(t.TypeParams, func(el types.TypeRef) any { return i.TypeRef(el) })
	}

	res["intro"] = i.TypeRef(t.Type)

	return res
}

func (i *Inspect) typeDefRef(t *types.TypeDefRef) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := i.typeRefs(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func (i *Inspect) TypeAliasRef(t *types.TypeAliasRef) map[string]any {
	d := t.Definition

	res := map[string]any{"typealias": Modifiers(t.Modifiers) + i.composeType(d.Package, d.TypeName, nil) + " = " + i.typeRef(d.Type)}

	res["intro"] = i.TypeRef(t.Type)

	return res
}

func (i *Inspect) typeAliasRef(t *types.TypeAliasRef) string {
	return Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
}

func (i *Inspect) MapRef(t *types.MapRef) map[string]any {
	res := map[string]any{
		"map":   Modifiers(t.Modifiers) + "map[..]..",
		"key":   i.TypeRef(t.Key),
		"value": i.TypeRef(t.Value),
	}

	return res
}

func (i *Inspect) mapRef(t *types.MapRef) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("map[%s]", i.typeRef(t.Key))
	out += fmt.Sprintf("%s", i.typeRef(t.Value))

	return out
}

func (i *Inspect) ChanRef(t *types.ChanRef) map[string]any {
	res := map[string]any{
		"chan": Modifiers(t.Modifiers) + "chan ..",
		"type": i.TypeRef(t.Value),
	}

	return res
}

func (i *Inspect) chanRef(t *types.ChanRef) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("chan %s", i.typeRef(t.Value))

	return out
}

func (i *Inspect) FuncTypeRef(t *types.FuncTypeRef) any {
	res := map[string]any{"func": Modifiers(t.Modifiers) + "{..}"}

	if len(t.Params) != 0 {
		res["params"] = i.Fields(t.Params)
	}

	if len(t.Results) != 0 {
		res["results"] = i.Fields(t.Results)
	}

	return res
}

func (i *Inspect) funcTypeRef(t *types.FuncTypeRef) string {
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

func (i *Inspect) StructTypeRef(t *types.StructTypeRef) any {
	res := map[string]any{"struct": Modifiers(t.Modifiers) + "{..}"}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) structTypeRef(t *types.StructTypeRef) string {
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

func (i *Inspect) InterfaceTypeRef(t *types.InterfaceTypeRef) any {
	res := map[string]any{"interface": Modifiers(t.Modifiers) + "{..}"}

	if len(t.Fields) != 0 {
		res["fields"] = i.Fields(t.Fields)
	}

	return res
}

func (i *Inspect) interfaceTypeRef(t *types.InterfaceTypeRef) string {
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

func (i *Inspect) typeParamRef(t *types.TypeParamRef) string {
	return fmt.Sprintf("%s%s", Modifiers(t.Modifiers), t.Name)
}

func (i *Inspect) typeRefs(t types.TypeRefs) []string {
	return slices.Map(t, func(el types.TypeRef) string { return i.typeRef(el) })
}
