package inspect

import (
	"fmt"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (i *Inspect) Function(f *types.Function) map[string]any {
	res := map[string]any{}

	if f.Receiver != nil {
		res["receiver"] = i.Fields([]*types.Field{f.Receiver})[0]
	}

	if len(f.Parameters) != 0 {
		res["params"] = i.Fields(f.Parameters)
	}

	if len(f.Results) != 0 {
		res["results"] = i.Fields(f.Results)
	}

	return map[string]any{"func " + i.funcID(f): res}
}

func (i *Inspect) funcID(f *types.Function) string {
	recv := ""
	if f.Receiver != nil {
		recv = "(" + i.typeRef(f.Receiver.Type) + ") "
	}

	if f.Package == i.trimPackage {
		return recv + f.FuncName
	}
	out := fmt.Sprintf("%s.%s", f.Package, f.FuncName)
	return recv + strings.TrimPrefix(out, i.trimPackage)
}
