package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	midPosition "github.com/sabahtalateh/toolboxgen/internal/mid/position"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"github.com/sabahtalateh/toolboxgen/internal/types/position"
	"go/ast"
)

func (c *Converter) TypeRef(ctx Context, expr ast.Expr) (types.TypeRef, error) {
	ref := mid.ParseTypeRef(ctx.Files(), expr)
	if err := ref.Error(); err != nil {
		return nil, err
	}

	return c.midTypeRef(ctx, ref)
}

func (c *Converter) midTypeRef(ctx Context, ref mid.TypeRef) (types.TypeRef, error) {
	switch r := ref.(type) {
	case *mid.Type:
		return c.midType(ctx.WithPosition(r.Position), r)
	case *mid.Map:
		return c.midMap(ctx.WithPosition(r.Position), r)
	case *mid.Chan:
		return c.midChan(ctx.WithPosition(r.Position), r)
	case *mid.FuncType:
		return c.midFuncType(ctx.WithPosition(r.Position), r)
	case *mid.StructType:
		return c.midStructType(ctx.WithPosition(r.Position), r)
	case *mid.InterfaceType:
		return c.midInterfaceType(ctx.WithPosition(r.Position), r)
	default:
		return nil, errors.Errorf(midPosition.OfTypeRef(ref), "unknown type %T", r)
	}
}

func (c *Converter) midType(ctx Context, midType *mid.Type) (types.TypeRef, error) {
	if midType.Package == "" {
		if def, ok := ctx.Defined(midType.TypeName); ok {
			return &types.TypeParamRef{
				Declared:  midType.Declared,
				Original:  midType.TypeName,
				Name:      def.Name,
				Modifiers: Modifiers(midType.Modifiers),
				Position:  midType.Position,
			}, nil
		}
	}

	typ, err := c.findType(ctx, ctx.ResolvePackage(midType.Package), midType.TypeName)
	if err != nil {
		return nil, err
	}

	switch t := typ.(type) {
	case *types.Builtin:
		return c.builtinRef(midType, t), nil
	case *types.Struct:
		return c.structRef(ctx, midType, t)
	case *types.Interface:
		return c.interfaceRef(ctx, midType, t)
	case *types.TypeDef:
		return c.typeDefRef(ctx, midType, t)
	case *types.TypeAlias:
		return c.typeAliasRef(midType, t), nil
	default:
		return nil, errors.Errorf(position.OfType(t), "unknown type %T", t)
	}
}

func (c *Converter) midMap(ctx Context, midType *mid.Map) (*types.MapRef, error) {
	var (
		res *types.MapRef
		err error
	)

	res = &types.MapRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
	}

	if res.Key, err = c.midTypeRef(ctx, midType.Key); err != nil {
		return nil, err
	}

	if res.Value, err = c.midTypeRef(ctx, midType.Value); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) midChan(ctx Context, midType *mid.Chan) (*types.ChanRef, error) {
	var (
		res *types.ChanRef
		err error
	)

	res = &types.ChanRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
	}

	if res.Value, err = c.midTypeRef(ctx, midType.Value); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) midFuncType(ctx Context, midType *mid.FuncType) (*types.FuncTypeRef, error) {
	var (
		res *types.FuncTypeRef
		err error
	)

	res = &types.FuncTypeRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
	}

	if res.Params, err = c.midFields(ctx, midType.Params...); err != nil {
		return nil, err
	}

	if res.Results, err = c.midFields(ctx, midType.Results...); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) midStructType(ctx Context, midType *mid.StructType) (*types.StructTypeRef, error) {
	var (
		res *types.StructTypeRef
		err error
	)

	res = &types.StructTypeRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
	}

	if res.Fields, err = c.midFields(ctx, midType.Fields...); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) midInterfaceType(ctx Context, midType *mid.InterfaceType) (*types.InterfaceTypeRef, error) {
	var (
		res *types.InterfaceTypeRef
		err error
	)

	res = &types.InterfaceTypeRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
	}

	if res.Fields, err = c.midFields(ctx, midType.Fields...); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) midFields(ctx Context, fields ...*mid.Field) (types.Fields, error) {
	var (
		res     types.Fields
		typeRef types.TypeRef
		err     error
	)

	for _, field := range fields {
		typeRef, err = c.midTypeRef(ctx, field.Type)
		if err != nil {
			return nil, err
		}
		res = append(res, &types.Field{
			Declared: field.Declared,
			Name:     field.Name,
			Type:     typeRef,
			Position: field.Position,
		})
	}

	return res, nil
}

func (c *Converter) builtinRef(midType *mid.Type, typ *types.Builtin) *types.BuiltinRef {
	return &types.BuiltinRef{
		Declared:   midType.Declared,
		Modifiers:  Modifiers(midType.Modifiers),
		TypeName:   typ.TypeName,
		Definition: typ,
		Position:   midType.Position,
	}
}

func (c *Converter) structRef(ctx Context, midType *mid.Type, typ *types.Struct) (*types.StructRef, error) {
	var (
		res *types.StructRef
		err error
	)

	res = &types.StructRef{
		Declared:   midType.Declared,
		Modifiers:  Modifiers(midType.Modifiers),
		Package:    typ.Package,
		TypeName:   typ.TypeName,
		Fields:     typ.Fields,
		Definition: typ,
		Position:   midType.Position,
	}

	if res.TypeParams, err = c.refTypeParams(ctx, typ.TypeParams, midType.TypeParams); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) interfaceRef(ctx Context, midType *mid.Type, typ *types.Interface) (*types.InterfaceRef, error) {
	var (
		res *types.InterfaceRef
		err error
	)

	res = &types.InterfaceRef{
		Declared:   midType.Declared,
		Modifiers:  Modifiers(midType.Modifiers),
		Package:    typ.Package,
		TypeName:   typ.TypeName,
		Definition: typ,
		Position:   midType.Position,
	}

	if res.TypeParams, err = c.refTypeParams(ctx, typ.TypeParams, midType.TypeParams); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) typeDefRef(ctx Context, midType *mid.Type, typ *types.TypeDef) (*types.TypeDefRef, error) {
	var (
		res *types.TypeDefRef
		err error
	)

	res = &types.TypeDefRef{
		Declared:   midType.Declared,
		Modifiers:  Modifiers(midType.Modifiers),
		Package:    typ.Package,
		TypeName:   typ.TypeName,
		Type:       typ.Type,
		Definition: typ,
		Position:   midType.Position,
	}

	if res.TypeParams, err = c.refTypeParams(ctx, typ.TypeParams, midType.TypeParams); err != nil {
		return nil, err
	}

	if err = forwardActual(ctx.WithDefined(res.Definition.TypeParams), res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) typeAliasRef(midType *mid.Type, typ *types.TypeAlias) *types.TypeAliasRef {
	return &types.TypeAliasRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Package:   typ.Package,
		TypeName:  typ.TypeName,
		Type:      typ.Type,
		Position:  midType.Position,
	}
}

func (c *Converter) refTypeParams(
	ctx Context,
	defined types.TypeParams,
	actual []mid.TypeRef,
) ([]types.TypeRef, error) {
	if len(defined) != len(actual) {
		return nil, errors.Errorf(ctx.Position(), "got %d type param but %d required", len(actual), len(defined))
	}

	var res []types.TypeRef
	for _, param := range actual {
		conv, err := c.midTypeRef(ctx, param)
		if err != nil {
			return nil, err
		}
		res = append(res, conv)
	}
	return res, nil
}

func forwardActual(ctx Context, outer *types.TypeDefRef) error {
	var (
		outerParams  = refTypeParams(outer)
		inner        = outer.Type
		innerParams  []types.TypeRef
		innerDefined types.TypeParams
	)

	switch in := inner.(type) {
	case *types.StructRef:
		innerParams = in.TypeParams
		innerDefined = in.Definition.TypeParams
	case *types.InterfaceRef:
		innerParams = in.TypeParams
		innerDefined = in.Definition.TypeParams
	case *types.TypeDefRef:
		innerParams = in.TypeParams
		innerDefined = in.Definition.TypeParams
	default:
		return nil
	}

	reverts := map[int]string{}
	for i, param := range innerParams {
		switch p := param.(type) {
		case *types.TypeParamRef:
			reverts[i] = innerDefined[i].Original
			innerDefined[i].Original = p.Original

			def, ok := ctx.Defined(p.Original)
			if !ok || len(outerParams)-1 < def.Order {
				return errors.Errorf(ctx.Position(), "inconsistent type params")
			}
			actual := outerParams[def.Order]
			innerParams[i] = actual
		}
	}

	switch in := inner.(type) {
	case *types.TypeDefRef:
		return forwardActual(ctx.WithDefined(innerDefined).WithPosition(position.OfTypeRef(inner)), in)
	}

	for i, orig := range reverts {
		innerDefined[i].Original = orig
	}

	return nil
}

func refTypeParams(ref types.TypeRef) []types.TypeRef {
	switch r := ref.(type) {
	case *types.StructRef:
		return r.TypeParams
	case *types.InterfaceRef:
		return r.TypeParams
	case *types.TypeDefRef:
		return r.TypeParams
	default:
		return nil
	}
}
