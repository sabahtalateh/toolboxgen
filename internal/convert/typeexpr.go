package convert

import (
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) TypeExpr(ctx Context, e ast.Expr) (types.TypeExpr, error) {
	typeExpr := syntax.ParseTypeExpr(ctx.Files(), e)
	if err := typeExpr.Error(); err != nil {
		return nil, err
	}

	return c.convertTypeExpr(ctx, typeExpr)
}

func (c *Converter) convertTypeExpr(ctx Context, e syntax.TypeExpr) (types.TypeExpr, error) {
	switch r := e.(type) {
	case *syntax.Type:
		return c.typeExpr(ctx.WithPosition(r.Position), r)
	case *syntax.Map:
		return c.mapExpr(ctx.WithPosition(r.Position), r)
	case *syntax.Chan:
		return c.chanExpr(ctx.WithPosition(r.Position), r)
	case *syntax.FuncType:
		return c.funcTypeExpr(ctx.WithPosition(r.Position), r)
	case *syntax.StructType:
		return c.structTypeExpr(ctx.WithPosition(r.Position), r)
	case *syntax.InterfaceType:
		return c.interfaceTypeExpr(ctx.WithPosition(r.Position), r)
	default:
		return nil, errors.Errorf(e.Get().Position(), "unknown type %T", r)
	}
}

func (c *Converter) typeExpr(ctx Context, e *syntax.Type) (types.TypeExpr, error) {
	if e.Package == "" {
		if def, ok := ctx.DefinedByName(e.TypeName); ok {
			return typeArgExpr(e, def), nil
		}
	}

	typ, err := c.findType(ctx, ctx.ResolvePackage(e.Package), e.TypeName)
	if err != nil {
		return nil, err
	}

	switch t := typ.(type) {
	case *types.Builtin:
		return builtinExpr(e, t), nil
	case *types.Struct:
		return c.structExpr(ctx, e, t)
	case *types.Interface:
		return c.interfaceExpr(ctx, e, t)
	case *types.TypeDef:
		return c.typeDefExpr(ctx, e, t)
	case *types.TypeAlias:
		return typeAliasExpr(e, t), nil
	default:
		return nil, errors.Errorf(t.Get().Position(), "unknown type %T", t)
	}
}

func (c *Converter) mapExpr(ctx Context, e *syntax.Map) (*types.MapExpr, error) {
	var (
		res *types.MapExpr
		err error
	)

	res = &types.MapExpr{
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

func (c *Converter) chanExpr(ctx Context, e *syntax.Chan) (*types.ChanExpr, error) {
	var (
		res *types.ChanExpr
		err error
	)

	res = &types.ChanExpr{
		Modifiers: Modifiers(e.Modifiers),
		Position:  e.Position,
		Code:      e.Code,
	}

	if res.Value, err = c.convertTypeExpr(ctx, e.Value); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) funcTypeExpr(ctx Context, e *syntax.FuncType) (*types.FuncTypeExpr, error) {
	var (
		res *types.FuncTypeExpr
		err error
	)

	res = &types.FuncTypeExpr{
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

func (c *Converter) structTypeExpr(ctx Context, e *syntax.StructType) (*types.StructTypeExpr, error) {
	var (
		res *types.StructTypeExpr
		err error
	)

	res = &types.StructTypeExpr{
		Modifiers: Modifiers(e.Modifiers),
		Position:  e.Position,
		Code:      e.Code,
	}

	if res.Fields, err = c.fields(ctx, e.Fields...); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) interfaceTypeExpr(ctx Context, e *syntax.InterfaceType) (*types.InterfaceTypeExpr, error) {
	var (
		res *types.InterfaceTypeExpr
		err error
	)

	res = &types.InterfaceTypeExpr{
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
		res      types.Fields
		typeExpr types.TypeExpr
		err      error
	)

	for _, field := range ff {
		if typeExpr, err = c.convertTypeExpr(ctx, field.Type); err != nil {
			return nil, err
		}
		res = append(res, &types.Field{
			Name:     field.Name,
			Type:     typeExpr,
			Position: field.Position,
			Code:     field.Code,
		})
	}

	return res, nil
}

func builtinExpr(e *syntax.Type, t *types.Builtin) *types.BuiltinExpr {
	return &types.BuiltinExpr{
		Modifiers:  Modifiers(e.Modifiers),
		TypeName:   t.TypeName,
		Position:   e.Position,
		Definition: t,
		Code:       e.Code,
	}
}

func (c *Converter) structExpr(ctx Context, e *syntax.Type, t *types.Struct) (*types.StructExpr, error) {
	actual, err := c.typeArgs(ctx, t.TypeParams, e.TypeArgs)
	if err != nil {
		return nil, err
	}

	return forwardToStruct(
		ctx.WithDefined(t.TypeParams),
		&types.StructExpr{
			Modifiers:  Modifiers(e.Modifiers),
			TypeArgs:   InitTypeArgs(t.TypeParams),
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

func (c *Converter) interfaceExpr(ctx Context, e *syntax.Type, t *types.Interface) (*types.InterfaceExpr, error) {
	actual, err := c.typeArgs(ctx, t.TypeParams, e.TypeArgs)
	if err != nil {
		return nil, err
	}

	return forwardToInterface(
		ctx.WithDefined(t.TypeParams),
		&types.InterfaceExpr{
			Modifiers:  Modifiers(e.Modifiers),
			TypeArgs:   InitTypeArgs(t.TypeParams),
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

func (c *Converter) typeDefExpr(ctx Context, e *syntax.Type, t *types.TypeDef) (*types.TypeDefExpr, error) {
	actual, err := c.typeArgs(ctx, t.TypeParams, e.TypeArgs)
	if err != nil {
		return nil, err
	}

	return forwardToTypeDef(
		ctx.WithDefined(t.TypeParams),
		&types.TypeDefExpr{
			Modifiers:  Modifiers(e.Modifiers),
			TypeArgs:   InitTypeArgs(t.TypeParams),
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

func typeAliasExpr(e *syntax.Type, t *types.TypeAlias) *types.TypeAliasExpr {
	return &types.TypeAliasExpr{
		Modifiers:  Modifiers(e.Modifiers),
		Package:    t.Package,
		TypeName:   t.TypeName,
		Type:       t.Type,
		Position:   e.Position,
		Definition: t,
		Code:       e.Code,
	}
}

func typeArgExpr(e *syntax.Type, t *types.TypeParam) *types.TypeArgExpr {
	return &types.TypeArgExpr{
		Modifiers:  Modifiers(e.Modifiers),
		Name:       t.Name,
		Order:      t.Order,
		Position:   e.Position,
		Definition: t,
		Code:       e.Code,
	}
}

func (c *Converter) typeArgs(ctx Context, defined types.TypeParams, ee syntax.TypeExprs) (types.TypeExprs, error) {
	if len(defined) != len(ee) {
		return nil, errors.Errorf(ctx.Position(), "got %d type param but %d required", len(ee), len(defined))
	}

	var res types.TypeExprs
	for _, param := range ee {
		act, err := c.convertTypeExpr(ctx, param)
		if err != nil {
			return nil, err
		}
		res = append(res, act)
	}
	return res, nil
}
