package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"go/ast"
)

func (c *Converter) TypeRef(ctx Context, expr ast.Expr) (types.TypeRef, error) {
	ref := mid.ParseTypeRef(ctx.Files, expr)
	if err := ref.ParseError(); err != nil {
		return nil, err
	}

	return c.convertMidTypeRef(ctx, ref)
}

func (c *Converter) convertMidTypeRef(ctx Context, ref mid.TypeRef) (types.TypeRef, error) {
	switch r := ref.(type) {
	case *mid.Type:
		return c.convertMidType(ctx, r)
	case *mid.Chan:
		println(ref)
	}

	return nil, nil
}

func (c *Converter) convertMidType(ctx Context, midType *mid.Type) (types.TypeRef, error) {
	if midType.Package == "" {
		if def, ok := ctx.Defined[midType.TypeName]; ok {
			return &types.TypeParamRef{
				Declared:  midType.Declared,
				Original:  midType.TypeName,
				Name:      def.Name,
				Modifiers: Modifiers(midType.Modifiers),
				Position:  midType.Position,
			}, nil
		}
	}

	typ, err := c.findType(ctx, ctx.PackageFromAlias(midType.Package), midType.TypeName)
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
		return nil, errors.Errorf(types.TypePosition(t), "unknown type %T", t)
	}
}

func (c *Converter) builtinRef(midType *mid.Type, typ *types.Builtin) *types.BuiltinRef {
	return &types.BuiltinRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		TypeName:  typ.TypeName,
		Position:  midType.Position,
	}
}

func (c *Converter) structRef(ctx Context, midType *mid.Type, typ *types.Struct) (*types.StructRef, error) {
	var (
		ref *types.StructRef
		err error
	)

	ref = &types.StructRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Package:   typ.Package,
		TypeName:  typ.TypeName,
		Position:  midType.Position,
	}

	ref.TypeParams, err = c.refTypeParams(ctx, typ.TypeParams, midType.TypeParams)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (c *Converter) interfaceRef(ctx Context, midType *mid.Type, typ *types.Interface) (*types.InterfaceRef, error) {
	var (
		ref *types.InterfaceRef
		err error
	)

	ref = &types.InterfaceRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Package:   typ.Package,
		TypeName:  typ.TypeName,
		Position:  midType.Position,
	}

	ref.TypeParams, err = c.refTypeParams(ctx, typ.TypeParams, midType.TypeParams)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (c *Converter) typeDefRef(ctx Context, midType *mid.Type, typ *types.TypeDef) (*types.TypeDefRef, error) {
	var (
		ref *types.TypeDefRef
		err error
	)

	ref = &types.TypeDefRef{
		Declared:  midType.Declared,
		Modifiers: Modifiers(midType.Modifiers),
		Package:   typ.Package,
		TypeName:  typ.TypeName,
		Type:      typ.Type,
		Position:  midType.Position,
	}

	ref.TypeParams, err = c.refTypeParams(ctx, typ.TypeParams, midType.TypeParams)
	if err != nil {
		return nil, err
	}

	forwardActual(ref)

	return ref, nil
}

func forwardActual(typeDef *types.TypeDefRef) {
	switch t := typeDef.Type.(type) {
	case *types.TypeDefRef:
		for i, param := range typeDef.TypeParams {
			t.TypeParams[i] = param
		}
		forwardActual(t)
	}
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
		return nil, errors.Errorf(ctx.CurrentPosition(), "got %d type param but %d required", len(actual), len(defined))
	}

	var res []types.TypeRef
	for _, param := range actual {
		conv, err := c.convertMidTypeRef(ctx, param)
		if err != nil {
			return nil, err
		}
		res = append(res, conv)
	}
	return res, nil
}
