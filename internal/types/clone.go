package types

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

func (t *Builtin) Clone() Type {
	return &Builtin{
		TypeName: t.TypeName,
		Declared: t.Declared,
	}
}

func (t *Struct) Clone() Type {
	return &Struct{
		TypeParams:   t.TypeParams.Clone(),
		Package:      t.Package,
		TypeName:     t.TypeName,
		Fields:       t.Fields.Clone(),
		Position:     t.Position,
		TypePosition: t.TypePosition,
		Declared:     t.Declared,
	}
}

func (t *Interface) Clone() Type {
	return &Interface{
		TypeParams:   t.TypeParams.Clone(),
		Package:      t.Package,
		TypeName:     t.TypeName,
		Methods:      t.Methods.Clone(),
		Position:     t.Position,
		TypePosition: t.TypePosition,
		Declared:     t.Declared,
	}
}

func (t *TypeDef) Clone() Type {
	return &TypeDef{
		TypeParams:   t.TypeParams.Clone(),
		Package:      t.Package,
		TypeName:     t.TypeName,
		Type:         t.Type.Clone(),
		Position:     t.Position,
		TypePosition: t.TypePosition,
		Declared:     t.Declared,
	}
}

func (t *TypeAlias) Clone() Type {
	return &TypeAlias{
		Package:      t.Package,
		TypeName:     t.TypeName,
		Type:         t.Type.Clone(),
		Position:     t.Position,
		TypePosition: t.TypePosition,
		Declared:     t.Declared,
	}
}

func (t *TypeParam) Clone() *TypeParam {
	return &TypeParam{
		Order:    t.Order,
		Name:     t.Name,
		Position: t.Position,
		Declared: t.Declared,
	}
}

func (t TypeParams) Clone() TypeParams {
	var res TypeParams
	for _, param := range t {
		res = append(res, param.Clone())
	}
	return res
}

func (t *BuiltinRef) Clone() TypeRef {
	def, ok := t.Definition.Clone().(*Builtin)
	if !ok {
		panic(errors.UnexpectedType(t.Position))
	}

	return &BuiltinRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeName:   t.TypeName,
		Position:   t.Position,
		Definition: def,
		Declared:   t.Declared,
	}
}

func (t *StructRef) Clone() TypeRef {
	def, ok := t.Definition.Clone().(*Struct)
	if !ok {
		panic(errors.UnexpectedType(t.Position))
	}

	return &StructRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Fields:     t.Fields.Clone(),
		Position:   t.Position,
		Definition: def,
		Declared:   t.Declared,
	}
}

func (t *InterfaceRef) Clone() TypeRef {
	def, ok := t.Definition.Clone().(*Interface)
	if !ok {
		panic(errors.UnexpectedType(t.Position))
	}

	return &InterfaceRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Methods:    t.Methods.Clone(),
		Position:   t.Position,
		Definition: def,
		Declared:   t.Declared,
	}
}

func (t *TypeDefRef) Clone() TypeRef {
	def, ok := t.Definition.Clone().(*TypeDef)
	if !ok {
		panic(errors.UnexpectedType(t.Position))
	}

	return &TypeDefRef{
		Modifiers:  t.Modifiers.Clone(),
		TypeParams: t.TypeParams.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       t.Type.Clone(),
		Position:   t.Position,
		Definition: def,
		Declared:   t.Declared,
	}
}

func (t *TypeAliasRef) Clone() TypeRef {
	def, ok := t.Definition.Clone().(*TypeAlias)
	if !ok {
		panic(errors.UnexpectedType(t.Position))
	}

	return &TypeAliasRef{
		Modifiers:  t.Modifiers.Clone(),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       t.Type.Clone(),
		Position:   t.Position,
		Definition: def,
		Declared:   t.Declared,
	}
}

func (t *MapRef) Clone() TypeRef {
	return &MapRef{
		Modifiers: t.Modifiers.Clone(),
		Key:       t.Key.Clone(),
		Value:     t.Value.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *ChanRef) Clone() TypeRef {
	return &ChanRef{
		Modifiers: t.Modifiers.Clone(),
		Value:     t.Value.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *FuncTypeRef) Clone() TypeRef {
	return &FuncTypeRef{
		Modifiers: t.Modifiers.Clone(),
		Params:    t.Params.Clone(),
		Results:   t.Results.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *StructTypeRef) Clone() TypeRef {
	return &StructTypeRef{
		Modifiers: t.Modifiers.Clone(),
		Fields:    t.Fields.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *InterfaceTypeRef) Clone() TypeRef {
	return &InterfaceTypeRef{
		Modifiers: t.Modifiers.Clone(),
		Fields:    t.Fields.Clone(),
		Position:  t.Position,
		Declared:  t.Declared,
	}
}

func (t *TypeParamRef) Clone() TypeRef {
	return &TypeParamRef{
		Modifiers:  t.Modifiers.Clone(),
		Order:      t.Order,
		Name:       t.Name,
		Position:   t.Position,
		Definition: t.Definition.Clone(),
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

func (f Fields) Clone() Fields {
	var res Fields
	for _, field := range f {
		res = append(res, field.Clone())
	}
	return res
}

func (f *Field) Clone() *Field {
	return &Field{
		Name:     f.Name,
		Type:     f.Type.Clone(),
		Position: f.Position,
		Declared: f.Declared,
	}
}
