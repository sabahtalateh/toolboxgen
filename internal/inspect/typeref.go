package inspect

import (
	"fmt"
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (i *Inspect) TypeRef(t types.TypeRef) string {
	switch tt := t.(type) {
	case *types.BuiltinRef:
		return i.BuiltinRef(tt)
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
		return i.TypeParamRef(tt)
	default:
		panic("unknown type reference")
	}
}

func (i *Inspect) BuiltinRef(t *types.BuiltinRef) string {
	return Modifiers(t.Modifiers) + t.TypeName
}

func (i *Inspect) StructRef(t *types.StructRef) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := i.TypeRefs(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func (i *Inspect) InterfaceRef(t *types.InterfaceRef) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := i.TypeRefs(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func (i *Inspect) TypeDefRef(t *types.TypeDefRef) string {
	out := Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) > 0 {
		typeParamsOut := i.TypeRefs(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	return out
}

func (i *Inspect) TypeAliasRef(t *types.TypeAliasRef) string {
	return Modifiers(t.Modifiers) + i.typeID(t.Package, t.TypeName)
}

func (i *Inspect) MapRef(t *types.MapRef) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("map[%s]", i.TypeRef(t.Key))
	out += fmt.Sprintf("%s", i.TypeRef(t.Value))

	return out
}

func (i *Inspect) ChanRef(t *types.ChanRef) string {
	out := Modifiers(t.Modifiers)
	out += fmt.Sprintf("chan %s", i.TypeRef(t.Value))

	return out
}

func (i *Inspect) FuncTypeRef(t *types.FuncTypeRef) string {
	out := Modifiers(t.Modifiers) + "func ("
	out += strings.Join(i.Fields2(t.Params), ", ") + ")"

	if len(t.Results) > 0 {
		results := i.Fields2(t.Results)
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

func (i *Inspect) StructTypeRef(t *types.StructTypeRef) string {
	out := Modifiers(t.Modifiers) + "struct{"

	if len(t.Fields) > 0 {
		fields := i.Fields2(t.Fields)
		for _, field := range fields {
			out += " " + field + ";"
		}
		out = strings.TrimSuffix(out, ";")
		out += " "
	}

	out += "}"

	return out
}

func (i *Inspect) InterfaceTypeRef(t *types.InterfaceTypeRef) string {
	out := Modifiers(t.Modifiers) + "interface{"

	if len(t.Fields) > 0 {
		fields := i.Fields2(t.Fields)
		for _, field := range fields {
			out += " " + field + ";"
		}
		out = strings.TrimSuffix(out, ";")
		out += " "
	}

	out += "}"

	return out
}

func (i *Inspect) TypeParamRef(t *types.TypeParamRef) string {
	return fmt.Sprintf("%s%s", Modifiers(t.Modifiers), t.Name)
}

func (i *Inspect) Fields2(f types.Fields) []string {
	return slices.Map(f, func(el *types.Field) string { return i.Field2(el) })
}

func (i *Inspect) TypeRefs(t types.TypeRefs) []string {
	return slices.Map(t, func(el types.TypeRef) string { return i.TypeRef(el) })
}
