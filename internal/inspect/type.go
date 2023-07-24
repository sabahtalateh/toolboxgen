package inspect

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"gopkg.in/yaml.v3"
)

func (i *Inspect) Type(t types.Type) string {
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

func (i *Inspect) Builtin(t *types.Builtin) string {
	return t.TypeName
}

const y = `
A:
  type: t
  type2: t2
  fields:
    - f1
    - f2
`

func (i *Inspect) Struct(t *types.Struct) string {
	x := &yaml.Node{}
	yaml.Unmarshal([]byte(y), x)
	println(x)

	Type := i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		Type += fmt.Sprintf("[%s]", strings.Join(i.TypeParams(t.TypeParams), ", "))
	}

	kk := kv(t.TypeName,
		kv("type", v(Type)),
		kv("type2", v(Type)),
		kv("fields", vv(i.Fields(t.Fields)...)),
	)
	println(kk)

	out := &yaml.Node{
		// Kind: yaml.MappingNode,
		// Content: kv(
		// 	t.TypeName,
		// 	kv("type", nn(Type)),
		// kv("type2", nn(Type)),
		// kv("fields", i.Fields(t.Fields)...),
		// ),
	}
	// Map := m(out, t.TypeName)
	// kv(Map, "type", Type)
	// kv(Map, "fields", i.Fields(t.Fields)...)

	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(out); err != nil {
		panic(err)
	}

	println(b.String())

	return ""
}

func (i *Inspect) Interface(t *types.Interface) string {
	out := i.typeBlock(t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		typeParamsOut := i.TypeParams(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	out += " interface {"
	methods := i.Fields2(t.Fields)

	if len(methods) > 0 {
		for _, method := range methods {
			out += "\n\t" + method
		}
		out += "\n"
	}

	out += "}"

	return out
}

func (i *Inspect) TypeDef(t *types.TypeDef) string {
	out := i.typeBlock(t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		typeParamsOut := i.TypeParams(t.TypeParams)
		out += fmt.Sprintf("[%s]", strings.Join(typeParamsOut, ", "))
	}
	out += " " + i.TypeRef(t.Type)

	return out
}

func (i *Inspect) TypeAlias(t *types.TypeAlias) string {
	out := i.typeBlock(t.Package, t.TypeName)
	out += " = "
	out += i.TypeRef(t.Type)

	return out
}

func (i *Inspect) TypeParam(t *types.TypeParam) string {
	return fmt.Sprintf("%s", t.Name)
}

func (i *Inspect) TypeParams(t types.TypeParams) []string {
	return slices.Map(t, func(el *types.TypeParam) string { return i.TypeParam(el) })
}
