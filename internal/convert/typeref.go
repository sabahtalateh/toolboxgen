package convert

import (
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	midPosition "github.com/sabahtalateh/toolboxgen/internal/mid/position"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) TypeRef(ctx Context, expr ast.Expr) (types.TypeRef, error) {
	midRef := mid.ParseTypeRef(ctx.Files(), expr)
	if err := midRef.Error(); err != nil {
		return nil, err
	}

	return c.midTypeRef(ctx, midRef)
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

func (c *Converter) midType(ctx Context, mid *mid.Type) (types.TypeRef, error) {
	if mid.Package == "" {
		if def, ok := ctx.DefinedByName(mid.TypeName); ok {
			return typeParamRef(mid, def), nil
		}
	}

	typ, err := c.findType(ctx, ctx.ResolvePackage(mid.Package), mid.TypeName)
	if err != nil {
		return nil, err
	}

	switch t := typ.(type) {
	case *types.Builtin:
		return builtinRef(mid, t), nil
	case *types.Struct:
		return c.structRef(ctx, mid, t)
	case *types.Interface:
		return c.interfaceRef(ctx, mid, t)
	case *types.TypeDef:
		return c.typeDefRef(ctx, mid, t)
	case *types.TypeAlias:
		return typeAliasRef(mid, t), nil
	default:
		return nil, errors.Errorf(t.Get().Position(), "unknown type %T", t)
	}
}

func (c *Converter) midMap(ctx Context, midType *mid.Map) (*types.MapRef, error) {
	var (
		res *types.MapRef
		err error
	)

	res = &types.MapRef{
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
		Declared:  midType.Declared,
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
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
		Declared:  midType.Declared,
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
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
		Declared:  midType.Declared,
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
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
		Declared:  midType.Declared,
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
		Modifiers: Modifiers(midType.Modifiers),
		Position:  midType.Position,
		Declared:  midType.Declared,
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
		if typeRef, err = c.midTypeRef(ctx, field.Type); err != nil {
			return nil, err
		}
		res = append(res, &types.Field{
			Name:     field.Name,
			Type:     typeRef,
			Position: field.Position,
			Declared: field.Declared,
		})
	}

	return res, nil
}

func builtinRef(midType *mid.Type, typ *types.Builtin) *types.BuiltinRef {
	return &types.BuiltinRef{
		Modifiers:  Modifiers(midType.Modifiers),
		TypeName:   typ.TypeName,
		Position:   midType.Position,
		Definition: typ,
		Declared:   midType.Declared,
	}
}

func (c *Converter) structRef(ctx Context, mid *mid.Type, typ *types.Struct) (*types.StructRef, error) {
	actual, err := c.actual(ctx, typ.TypeParams, mid.TypeParams)
	if err != nil {
		return nil, err
	}

	return resolveStruct(
		ctx.WithDefined(typ.TypeParams),
		&types.StructRef{
			Modifiers:  Modifiers(mid.Modifiers),
			TypeParams: InitTypeParams(typ.TypeParams),
			Package:    typ.Package,
			TypeName:   typ.TypeName,
			Fields:     typ.Fields.Clone(),
			Position:   mid.Position,
			Definition: typ,
			Declared:   mid.Declared,
		},
		actual,
	)
}

func (c *Converter) interfaceRef(ctx Context, mid *mid.Type, typ *types.Interface) (*types.InterfaceRef, error) {
	actual, err := c.actual(ctx, typ.TypeParams, mid.TypeParams)
	if err != nil {
		return nil, err
	}

	return resolveInterface(
		ctx.WithDefined(typ.TypeParams),
		&types.InterfaceRef{
			Modifiers:  Modifiers(mid.Modifiers),
			TypeParams: InitTypeParams(typ.TypeParams),
			Package:    typ.Package,
			TypeName:   typ.TypeName,
			Methods:    typ.Methods.Clone(),
			Position:   mid.Position,
			Definition: typ,
			Declared:   mid.Declared,
		},
		actual,
	)
}

func (c *Converter) typeDefRef(ctx Context, mid *mid.Type, typ *types.TypeDef) (*types.TypeDefRef, error) {
	actual, err := c.actual(ctx, typ.TypeParams, mid.TypeParams)
	if err != nil {
		return nil, err
	}

	return resolveTypeDef(
		ctx.WithDefined(typ.TypeParams),
		&types.TypeDefRef{
			Modifiers:  Modifiers(mid.Modifiers),
			TypeParams: InitTypeParams(typ.TypeParams),
			Package:    typ.Package,
			TypeName:   typ.TypeName,
			Type:       typ.Type.Clone(),
			Position:   mid.Position,
			Definition: typ,
			Declared:   mid.Declared,
		},
		actual,
	)
}

func typeAliasRef(mid *mid.Type, typ *types.TypeAlias) *types.TypeAliasRef {
	return &types.TypeAliasRef{
		Modifiers:  Modifiers(mid.Modifiers),
		Package:    typ.Package,
		TypeName:   typ.TypeName,
		Type:       typ.Type,
		Position:   mid.Position,
		Definition: typ,
		Declared:   mid.Declared,
	}
}

func typeParamRef(mid *mid.Type, typ *types.TypeParam) *types.TypeParamRef {
	return &types.TypeParamRef{
		Modifiers:  Modifiers(mid.Modifiers),
		Name:       typ.Name,
		Order:      typ.Order,
		Position:   mid.Position,
		Definition: typ,
		Declared:   mid.Declared,
	}
}

func (c *Converter) actual(ctx Context, defined types.TypeParams, mids mid.TypeRefs) (types.TypeRefs, error) {
	if len(defined) != len(mids) {
		return nil, errors.Errorf(ctx.Position(), "got %d type param but %d required", len(mids), len(defined))
	}

	var res types.TypeRefs
	for _, param := range mids {
		act, err := c.midTypeRef(ctx, param)
		if err != nil {
			return nil, err
		}
		res = append(res, act)
	}
	return res, nil
}
