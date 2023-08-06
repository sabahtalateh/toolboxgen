package inspect

import (
	"fmt"
	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (i *Inspect) Fields(ff types.Fields) []any {
	return slices.Map(ff, func(f *types.Field) any {
		intro := i.TypeExpr(f.Type)
		switch in := intro.(type) {
		case map[string]any:
			res := map[string]any{}
			if f.Name != "" {
				res["name"] = f.Name
			}
			for k, v := range in {
				res[k] = v
			}

			return res
		default:
			res := ""
			if f.Name != "" {
				res += f.Name + " "
			}
			res += fmt.Sprintf("%s", in)

			return res
		}
	})
}

func (i *Inspect) field(f *types.Field) string {
	name := ""
	if f.Name != "" {
		name = f.Name + " "
	}
	return name + i.typeExpr(f.Type)
}

func (i *Inspect) fields(f types.Fields) []string {
	return slices.Map(f, func(el *types.Field) string { return i.field(el) })
}
