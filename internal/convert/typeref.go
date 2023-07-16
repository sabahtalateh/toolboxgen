package convert

import (
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"github.com/sabahtalateh/toolboxgen/internal/types/position"
)

func (c *Converter) TypeRef(ctx Context, expr ast.Expr) (types.TypeRef, error) {
	ref := mid.ParseTypeRef(ctx.Files(), expr)
	if err := ref.ParseError(); err != nil {
		return nil, err
	}

	return c.convertMidTypeRef(ctx, ref)
}

func (c *Converter) convertMidTypeRef(ctx Context, ref mid.TypeRef) (types.TypeRef, error) {
	switch r := ref.(type) {
	case *mid.Type:
		return c.convertMidType(ctx.WithPosition(r.Position), r)
	case *mid.Chan:
		println(ref)
	}

	return nil, nil
}

func (c *Converter) convertMidType(ctx Context, midType *mid.Type) (types.TypeRef, error) {
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
		ref *types.StructRef
		err error
	)

	ref = &types.StructRef{
		Declared:   midType.Declared,
		Modifiers:  Modifiers(midType.Modifiers),
		Package:    typ.Package,
		TypeName:   typ.TypeName,
		Definition: typ,
		Position:   midType.Position,
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
		Declared:   midType.Declared,
		Modifiers:  Modifiers(midType.Modifiers),
		Package:    typ.Package,
		TypeName:   typ.TypeName,
		Definition: typ,
		Position:   midType.Position,
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
		Declared:   midType.Declared,
		Modifiers:  Modifiers(midType.Modifiers),
		Package:    typ.Package,
		TypeName:   typ.TypeName,
		Type:       typ.Type,
		Definition: typ,
		Position:   midType.Position,
	}

	ref.TypeParams, err = c.refTypeParams(ctx, typ.TypeParams, midType.TypeParams)
	if err != nil {
		return nil, err
	}

	if err = forwardActual(ctx.WithDefined(ref.Definition.TypeParams), ref); err != nil {
		return nil, err
	}

	return ref, nil
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
		conv, err := c.convertMidTypeRef(ctx, param)
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
