package inspect

import (
	"fmt"
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (i *Inspect) Type(t types.Type) Out {
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

func (i *Inspect) Builtin(t *types.Builtin) Out {
	return Out{Yaml: kv(t.TypeName, kv("builtin", v(t.TypeName))).yaml()}
}

func (i *Inspect) Struct(t *types.Struct) Out {
	Type := i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		Type += fmt.Sprintf("[%s]", strings.Join(i.TypeParams(t.TypeParams), ", "))
	}

	out := kv(t.TypeName,
		kv("type", v(Type)),
		// kv("fields", vv(i.Fields(t.Fields)...)),
		kv("fields", vv(
			kv("method", kv("actual", vv("A", "B"))),
		)),
	)

	return Out{Yaml: out.yaml()}
}

func (i *Inspect) Interface(t *types.Interface) Out {
	Type := i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		Type += fmt.Sprintf("[%s]", strings.Join(i.TypeParams(t.TypeParams), ", "))
	}

	out := kv(t.TypeName,
		kv("interface", v(Type)),
		kv("fields", vv(i.Fields(t.Fields)...)),
	)

	return Out{Yaml: out.yaml()}
}

func (i *Inspect) TypeDef(t *types.TypeDef) Out {
	Type := i.typeID(t.Package, t.TypeName)
	if len(t.TypeParams) != 0 {
		Type += fmt.Sprintf("[%s]", strings.Join(i.TypeParams(t.TypeParams), ", "))
	}

	out := kv(t.TypeName,
		kv("typedef", v(Type)),
		kv("type", v(i.TypeRef(t.Type))),
	)

	return Out{Yaml: out.yaml()}
}

func (i *Inspect) TypeAlias(t *types.TypeAlias) Out {
	out := kv(t.TypeName,
		kv("typealias", v(i.typeID(t.Package, t.TypeName))),
		kv("type", v(i.TypeRef(t.Type))),
	)

	return Out{Yaml: out.yaml()}
}

func (i *Inspect) TypeParam(t *types.TypeParam) string {
	return fmt.Sprintf("%s", t.Name)
}

func (i *Inspect) TypeParams(t types.TypeParams) []string {
	return slices.Map(t, func(el *types.TypeParam) string { return i.TypeParam(el) })
}
