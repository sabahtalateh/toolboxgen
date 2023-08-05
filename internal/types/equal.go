package types

import "github.com/sabahtalateh/toolboxgen/internal/maps"

func (t *Builtin) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *Builtin:
		return t.TypeName == tt2.TypeName
	default:
		return false
	}
}

func (t *Struct) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *Struct:
		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		if !t.Fields.Equal(tt2.Fields) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *Interface) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *Interface:
		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		methods1 := maps.FromSlice(t.Fields, func(s *Field) (string, *Field) { return s.Name, s })
		methods2 := maps.FromSlice(tt2.Fields, func(s *Field) (string, *Field) { return s.Name, s })

		if len(methods1) != len(methods2) {
			return false
		}

		for k, v := range methods1 {
			v2, ok := methods2[k]
			if !ok {
				return false
			}
			if !v.Equal(v2) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func (t *TypeDef) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *TypeDef:
		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		if !t.Type.Equal(tt2.Type) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeAlias) Equal(t2 Type) bool {
	switch tt2 := t2.(type) {
	case *TypeAlias:
		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.Type.Equal(tt2.Type) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeParam) Equal(t2 *TypeParam) bool {
	return t.Order == t2.Order
}

func (t TypeParams) Equal(t2 TypeParams) bool {
	if len(t) != len(t2) {
		return false
	}

	for i, param := range t {
		if !param.Equal(t2[i]) {
			return false
		}
	}

	return true
}

func (t *BuiltinExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *BuiltinExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *StructExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *StructExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *InterfaceExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *InterfaceExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeDefExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *TypeDefExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		if !t.TypeParams.Equal(tt2.TypeParams) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeAliasExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *TypeAliasExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Package != tt2.Package {
			return false
		}

		if t.TypeName != tt2.TypeName {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *MapExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *MapExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Key.Equal(tt2.Key) {
			return false
		}

		if !t.Value.Equal(tt2.Value) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *ChanExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *ChanExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Value.Equal(tt2.Value) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *FuncTypeExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *FuncTypeExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Params.Equal(tt2.Params) {
			return false
		}

		if !t.Results.Equal(tt2.Results) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *StructTypeExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *StructTypeExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Fields.Equal(tt2.Fields) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *InterfaceTypeExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *InterfaceTypeExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if !t.Fields.Equal(tt2.Fields) {
			return false
		}

		return true
	default:
		return false
	}
}

func (t *TypeParamExpr) Equal(t2 TypeExpr) bool {
	switch tt2 := t2.(type) {
	case *TypeParamExpr:
		if !t.Modifiers.Equal(tt2.Modifiers) {
			return false
		}

		if t.Order != tt2.Order {
			return false
		}

		return true
	default:
		return false
	}
}

func (t TypeExprs) Equal(t2 TypeExprs) bool {
	if len(t) != len(t2) {
		return false
	}

	for i, expr := range t {
		if !expr.Equal(t2[i]) {
			return false
		}
	}

	return true
}
