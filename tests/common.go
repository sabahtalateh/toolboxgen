package tests

import (
	"fmt"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func typeID(t types.Type) string {
	switch t := t.(type) {
	case *types.Builtin:
		return t.TypeName
	case *types.Struct:
		return fmt.Sprintf("%s.%s", t.Package, t.TypeName)
	case *types.Interface:
		return fmt.Sprintf("%s.%s", t.Package, t.TypeName)
	case *types.TypeDef:
		return fmt.Sprintf("%s.%s", t.Package, t.TypeName)
	case *types.TypeAlias:
		return fmt.Sprintf("%s.%s", t.Package, t.TypeName)
	default:
		panic("unknown type")
	}
}

type convertOut struct {
	types map[string]types.Type
	funcs map[string]*types.Function
}

func (c convertOut) equal(c2 convertOut) bool {
	for k2, v2 := range c2.types {
		v1, ok := c.types[k2]
		if !ok {
			return false
		}
		if !v1.Equal(v2) {
			return false
		}
	}

	return true
}
