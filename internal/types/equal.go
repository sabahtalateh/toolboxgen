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

func (t *BuiltinRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *BuiltinRef:
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

func (t *StructRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *StructRef:
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

func (t *InterfaceRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *InterfaceRef:
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

func (t *TypeDefRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *TypeDefRef:
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

func (t *TypeAliasRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *TypeAliasRef:
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

func (t *MapRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *MapRef:
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

func (t *ChanRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *ChanRef:
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

func (t *FuncTypeRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *FuncTypeRef:
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

func (t *StructTypeRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *StructTypeRef:
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

func (t *InterfaceTypeRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *InterfaceTypeRef:
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

func (t *TypeParamRef) Equal(t2 TypeRef) bool {
	switch tt2 := t2.(type) {
	case *TypeParamRef:
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

func (t TypeRefs) Equal(t2 TypeRefs) bool {
	if len(t) != len(t2) {
		return false
	}

	for i, ref := range t {
		if !ref.Equal(t2[i]) {
			return false
		}
	}

	return true
}
