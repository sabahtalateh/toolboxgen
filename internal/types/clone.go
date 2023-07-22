package types

func (t *BuiltinRef) Clone() TypeRef {
	if t == nil {
		return nil
	}

	return &BuiltinRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeName:   t.TypeName,
		Position:   t.Position,
		Definition: t.Definition,
		Declared:   t.Declared,
	}
}

func (t *StructRef) Clone() TypeRef {
	if t == nil {
		return nil
	}

	return &StructRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Fields:     t.Fields.Clone(),
		Position:   t.Position,
		Definition: t.Definition,
		Declared:   t.Declared,
	}
}

func (t *InterfaceRef) Clone() TypeRef {
	if t == nil {
		return nil
	}

	return &InterfaceRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Methods:    t.Methods.Clone(),
		Position:   t.Position,
		Definition: t.Definition,
		Declared:   t.Declared,
	}
}

func (t *TypeDefRef) Clone() TypeRef {
	if t == nil {
		return nil
	}

	return &TypeDefRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       t.Type.Clone(),
		Position:   t.Position,
		Definition: t.Definition,
		Declared:   t.Declared,
	}
}

func (t *TypeAliasRef) Clone() TypeRef {
	if t == nil {
		return nil
	}

	return &TypeAliasRef{
		Modifiers:  t.Modifiers.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       t.Type.Clone(),
		Position:   t.Position,
		Definition: t.Definition,
		Declared:   t.Declared,
	}
}

func (t *MapRef) Clone() TypeRef {
	if t == nil {
		return nil
	}
	return &MapRef{
		Modifiers: t.Modifiers.Clone(),
		Key:       t.Key.Clone(),
		Value:     t.Value.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *ChanRef) Clone() TypeRef {
	if t == nil {
		return nil
	}
	return &ChanRef{
		Modifiers: t.Modifiers.Clone(),
		Value:     t.Value.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *FuncTypeRef) Clone() TypeRef {
	if t == nil {
		return nil
	}
	return &FuncTypeRef{
		Modifiers: t.Modifiers.Clone(),
		Params:    t.Params.Clone(),
		Results:   t.Results.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *StructTypeRef) Clone() TypeRef {
	if t == nil {
		return nil
	}
	return &StructTypeRef{
		Modifiers: t.Modifiers.Clone(),
		Fields:    t.Fields.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *InterfaceTypeRef) Clone() TypeRef {
	if t == nil {
		return nil
	}
	return &InterfaceTypeRef{
		Modifiers: t.Modifiers.Clone(),
		Fields:    t.Fields.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *TypeParamRef) Clone() TypeRef {
	if t == nil {
		return nil
	}
	return &TypeParamRef{
		Modifiers:  t.Modifiers.Clone(),
		Order:      t.Order,
		Name:       t.Name,
		Position:   t.Position,
		Definition: t.Definition,
		Declared:   t.Declared,
	}
}

func (t TypeRefs) Clone() TypeRefs {
	var res TypeRefs
	for _, ref := range t {
		res = append(res, ref.Clone())
	}
	return res
}

func (m *Pointer) Clone() Modifier {
	if m == nil {
		return nil
	}
	return &Pointer{Position: m.Position}
}

func (m *Array) Clone() Modifier {
	if m == nil {
		return nil
	}
	return &Array{Sized: m.Sized, Position: m.Position}
}

func (m *Ellipsis) Clone() Modifier {
	if m == nil {
		return nil
	}
	return &Ellipsis{Position: m.Position}
}

func (m Modifiers) Clone() Modifiers {
	var res Modifiers
	for _, mod := range m {
		res = append(res, mod.Clone())
	}
	return res
}

func (f *Field) Clone() *Field {
	if f == nil {
		return nil
	}
	return &Field{
		Name:     f.Name,
		Type:     f.Type.Clone(),
		Position: f.Position,
		Declared: f.Declared,
	}
}

func (f Fields) Clone() Fields {
	var res Fields
	for _, field := range f {
		res = append(res, field.Clone())
	}
	return res
}
