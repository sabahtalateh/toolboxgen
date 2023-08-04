package convert

import (
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) TypeRef(ctx Context, e ast.Expr) (types.TypeRef, error) {
	typeExpr := syntax.ParseTypeExpr(ctx.Files(), e)
	if err := typeExpr.Error(); err != nil {
		return nil, err
	}

	return c.convertTypeExpr(ctx, typeExpr)
}

func (c *Converter) convertTypeExpr(ctx Context, e syntax.TypeExpr) (types.TypeRef, error) {
	switch r := e.(type) {
	case *syntax.Type:
		return c.typeRef(ctx.WithPosition(r.Position), r)
	case *syntax.Map:
		return c.mapRef(ctx.WithPosition(r.Position), r)
	case *syntax.Chan:
		return c.chanRef(ctx.WithPosition(r.Position), r)
	case *syntax.FuncType:
		return c.funcTypeRef(ctx.WithPosition(r.Position), r)
	case *syntax.StructType:
		return c.structTypeRef(ctx.WithPosition(r.Position), r)
	case *syntax.InterfaceType:
		return c.interfaceTypeRef(ctx.WithPosition(r.Position), r)
	default:
		return nil, errors.Errorf(e.Get().Position(), "unknown type %T", r)
	}
}

func (c *Converter) typeRef(ctx Context, e *syntax.Type) (types.TypeRef, error) {
	if e.Package == "" {
		if def, ok := ctx.DefinedByName(e.TypeName); ok {
			return typeParamRef(e, def), nil
		}
	}

	typ, err := c.findType(ctx, ctx.ResolvePackage(e.Package), e.TypeName)
	if err != nil {
		return nil, err
	}

	switch t := typ.(type) {
	case *types.Builtin:
		return builtinRef(e, t), nil
	case *types.Struct:
		return c.structRef(ctx, e, t)
	case *types.Interface:
		return c.interfaceRef(ctx, e, t)
	case *types.TypeDef:
		return c.typeDefRef(ctx, e, t)
	case *types.TypeAlias:
		return typeAliasRef(e, t), nil
	default:
		return nil, errors.Errorf(t.Get().Position(), "unknown type %T", t)
	}
}

func (c *Converter) mapRef(ctx Context, e *syntax.Map) (*types.MapRef, error) {
	var (
		res *types.MapRef
		err error
	)

	res = &types.MapRef{
		Modifiers: Modifiers(e.Modifiers),
		Position:  e.Position,
		Code:      e.Code,
	}

	if res.Key, err = c.convertTypeExpr(ctx, e.Key); err != nil {
		return nil, err
	}

	if res.Value, err = c.convertTypeExpr(ctx, e.Value); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) chanRef(ctx Context, e *syntax.Chan) (*types.ChanRef, error) {
	var (
		res *types.ChanRef
		err error
	)

	res = &types.ChanRef{
		Modifiers: Modifiers(e.Modifiers),
		Position:  e.Position,
		Code:      e.Code,
	}

	if res.Value, err = c.convertTypeExpr(ctx, e.Value); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) funcTypeRef(ctx Context, e *syntax.FuncType) (*types.FuncTypeRef, error) {
	var (
		res *types.FuncTypeRef
		err error
	)

	res = &types.FuncTypeRef{
		Modifiers: Modifiers(e.Modifiers),
		Position:  e.Position,
		Code:      e.Code,
	}

	if res.Params, err = c.fields(ctx, e.Params...); err != nil {
		return nil, err
	}

	if res.Results, err = c.fields(ctx, e.Results...); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) structTypeRef(ctx Context, e *syntax.StructType) (*types.StructTypeRef, error) {
	var (
		res *types.StructTypeRef
		err error
	)

	res = &types.StructTypeRef{
		Modifiers: Modifiers(e.Modifiers),
		Position:  e.Position,
		Code:      e.Code,
	}

	if res.Fields, err = c.fields(ctx, e.Fields...); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) interfaceTypeRef(ctx Context, e *syntax.InterfaceType) (*types.InterfaceTypeRef, error) {
	var (
		res *types.InterfaceTypeRef
		err error
	)

	res = &types.InterfaceTypeRef{
		Modifiers: Modifiers(e.Modifiers),
		Position:  e.Position,
		Code:      e.Code,
	}

	if res.Fields, err = c.fields(ctx, e.Fields...); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) fields(ctx Context, ff ...*syntax.Field) (types.Fields, error) {
	var (
		res     types.Fields
		typeRef types.TypeRef
		err     error
	)

	for _, field := range ff {
		if typeRef, err = c.convertTypeExpr(ctx, field.Type); err != nil {
			return nil, err
		}
		res = append(res, &types.Field{
			Name:     field.Name,
			Type:     typeRef,
			Position: field.Position,
			Code:     field.Code,
		})
	}

	return res, nil
}

func builtinRef(e *syntax.Type, t *types.Builtin) *types.BuiltinRef {
	return &types.BuiltinRef{
		Modifiers:  Modifiers(e.Modifiers),
		TypeName:   t.TypeName,
		Position:   e.Position,
		Definition: t,
		Code:       e.Code,
	}
}

func (c *Converter) structRef(ctx Context, e *syntax.Type, t *types.Struct) (*types.StructRef, error) {
	actual, err := c.actual(ctx, t.TypeParams, e.TypeParams)
	if err != nil {
		return nil, err
	}

	return resolveStruct(
		ctx.WithDefined(t.TypeParams),
		&types.StructRef{
			Modifiers:  Modifiers(e.Modifiers),
			TypeParams: InitTypeParams(t.TypeParams),
			Package:    t.Package,
			TypeName:   t.TypeName,
			Fields:     t.Fields.Clone(),
			Position:   e.Position,
			Definition: t,
			Code:       e.Code,
		},
		actual,
	)
}

func (c *Converter) interfaceRef(ctx Context, e *syntax.Type, t *types.Interface) (*types.InterfaceRef, error) {
	actual, err := c.actual(ctx, t.TypeParams, e.TypeParams)
	if err != nil {
		return nil, err
	}

	return resolveInterface(
		ctx.WithDefined(t.TypeParams),
		&types.InterfaceRef{
			Modifiers:  Modifiers(e.Modifiers),
			TypeParams: InitTypeParams(t.TypeParams),
			Package:    t.Package,
			TypeName:   t.TypeName,
			Fields:     t.Fields.Clone(),
			Position:   e.Position,
			Definition: t,
			Code:       e.Code,
		},
		actual,
	)
}

func (c *Converter) typeDefRef(ctx Context, e *syntax.Type, t *types.TypeDef) (*types.TypeDefRef, error) {
	actual, err := c.actual(ctx, t.TypeParams, e.TypeParams)
	if err != nil {
		return nil, err
	}

	return resolveTypeDef(
		ctx.WithDefined(t.TypeParams),
		&types.TypeDefRef{
			Modifiers:  Modifiers(e.Modifiers),
			TypeParams: InitTypeParams(t.TypeParams),
			Package:    t.Package,
			TypeName:   t.TypeName,
			Type:       t.Type.Clone(),
			Position:   e.Position,
			Definition: t,
			Code:       e.Code,
		},
		actual,
	)
}

func typeAliasRef(e *syntax.Type, t *types.TypeAlias) *types.TypeAliasRef {
	return &types.TypeAliasRef{
		Modifiers:  Modifiers(e.Modifiers),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       t.Type,
		Position:   e.Position,
		Definition: t,
		Code:       e.Code,
	}
}

func typeParamRef(e *syntax.Type, t *types.TypeParam) *types.TypeParamRef {
	return &types.TypeParamRef{
		Modifiers:  Modifiers(e.Modifiers),
		Name:       t.Name,
		Order:      t.Order,
		Position:   e.Position,
		Definition: t,
		Code:       e.Code,
	}
}

func (c *Converter) actual(ctx Context, defined types.TypeParams, ee syntax.TypeRefs) (types.TypeRefs, error) {
	if len(defined) != len(ee) {
		return nil, errors.Errorf(ctx.Position(), "got %d type param but %d required", len(ee), len(defined))
	}

	var res types.TypeRefs
	for _, param := range ee {
		act, err := c.convertTypeExpr(ctx, param)
		if err != nil {
			return nil, err
		}
		res = append(res, act)
	}
	return res, nil
}
