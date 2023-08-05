package types

func (t *BuiltinExpr) Clone() TypeExpr {
	return &BuiltinExpr{
		Modifiers:  t.Modifiers.Clone(),
		TypeName:   t.TypeName,
		Position:   t.Position,
		Definition: t.Definition,
		Code:       t.Code,
	}
}

func (t *StructExpr) Clone() TypeExpr {
	return &StructExpr{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Fields:     t.Fields.Clone(),
		Position:   t.Position,
		Definition: t.Definition,
		Code:       t.Code,
	}
}

func (t *InterfaceExpr) Clone() TypeExpr {
	return &InterfaceExpr{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Fields:     t.Fields.Clone(),
		Position:   t.Position,
		Definition: t.Definition,
		Code:       t.Code,
	}
}

func (t *TypeDefExpr) Clone() TypeExpr {
	var typ TypeExpr
	if t.Type != nil {
		typ = t.Type.Clone()
	}

	return &TypeDefExpr{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       typ,
		Position:   t.Position,
		Definition: t.Definition,
		Code:       t.Code,
	}
}

func (t *TypeAliasExpr) Clone() TypeExpr {
	var typ TypeExpr
	if t.Type != nil {
		typ = t.Type.Clone()
	}

	return &TypeAliasExpr{
		Modifiers:  t.Modifiers.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       typ,
		Position:   t.Position,
		Definition: t.Definition,
		Code:       t.Code,
	}
}

func (t *MapExpr) Clone() TypeExpr {
	return &MapExpr{
		Modifiers: t.Modifiers.Clone(),
		Key:       t.Key.Clone(),
		Value:     t.Value.Clone(),
		Position:  t.Position,
		Code:      t.Code,
	}
}

func (t *ChanExpr) Clone() TypeExpr {
	return &ChanExpr{
		Modifiers: t.Modifiers.Clone(),
		Value:     t.Value.Clone(),
		Position:  t.Position,
		Code:      t.Code,
	}
}

func (t *FuncTypeExpr) Clone() TypeExpr {
	return &FuncTypeExpr{
		Modifiers: t.Modifiers.Clone(),
		Params:    t.Params.Clone(),
		Results:   t.Results.Clone(),
		Position:  t.Position,
		Code:      t.Code,
	}
}

func (t *StructTypeExpr) Clone() TypeExpr {
	return &StructTypeExpr{
		Modifiers: t.Modifiers.Clone(),
		Fields:    t.Fields.Clone(),
		Position:  t.Position,
		Code:      t.Code,
	}
}

func (t *InterfaceTypeExpr) Clone() TypeExpr {
	return &InterfaceTypeExpr{
		Modifiers: t.Modifiers.Clone(),
		Fields:    t.Fields.Clone(),
		Position:  t.Position,
		Code:      t.Code,
	}
}

func (t *TypeParamExpr) Clone() TypeExpr {
	return &TypeParamExpr{
		Modifiers:  t.Modifiers.Clone(),
		Order:      t.Order,
		Name:       t.Name,
		Position:   t.Position,
		Definition: t.Definition,
		Code:       t.Code,
	}
}

func (t TypeExprs) Clone() TypeExprs {
	var res TypeExprs
	for _, expr := range t {
		res = append(res, expr.Clone())
	}
	return res
}

func (m *Pointer) Clone() Modifier {
	return &Pointer{Position: m.Position}
}

func (m *Array) Clone() Modifier {
	return &Array{Sized: m.Sized, Position: m.Position}
}

func (m *Ellipsis) Clone() Modifier {
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
	return &Field{
		Name:     f.Name,
		Type:     f.Type.Clone(),
		Position: f.Position,
		Code:     f.Code,
	}
}

func (f Fields) Clone() Fields {
	var res Fields
	for _, field := range f {
		res = append(res, field.Clone())
	}
	return res
}
