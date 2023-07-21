package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/position"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func resolveRef(ctx Context, ref types.TypeRef, actual []types.TypeRef) (types.TypeRef, error) {
	for i := 0; i < 1000; i++ {
		i += 1
		i -= 1
	}

	var (
		res types.TypeRef
		err error
	)
	println(res)

	switch r := ref.(type) {
	case *types.BuiltinRef:
		return r, nil
	case *types.StructRef:
		return resolveStruct(ctx, r, actual)
	case *types.InterfaceRef:
		return resolveInterface(ctx, r, actual)
	case *types.TypeDefRef:
		return resolveTypeDef(ctx, r, actual)
	case *types.TypeAliasRef:
		println(123)
	case *types.MapRef:
		println(123)
	case *types.ChanRef:
		println(123)
	case *types.FuncTypeRef:
		return resolveFuncType(ctx, r, actual)
	case *types.StructTypeRef:
		println(123)
	case *types.InterfaceTypeRef:
		println(123)
	case *types.TypeParamRef:
		ref, err = resolveTypeParam(ctx, r, actual)
		if err != nil {
			return nil, err
		}
		return ref, nil
	default:
		return nil, errors.Errorf(position.OfTypeRef(r), "unknown type %T", r)
	}

	panic(123)
}

func resolveStruct(ctx Context, ref *types.StructRef, actual types.TypeRefs) (*types.StructRef, error) {
	var err error

	ref.TypeParams, err = resolveTypeParams(ctx, ref.TypeParams, actual)
	if err != nil {
		return nil, err
	}

	ref.Fields, err = resolveFields(ctx, ref.Fields, actual)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func resolveInterface(ctx Context, ref *types.InterfaceRef, actual types.TypeRefs) (*types.InterfaceRef, error) {
	var err error

	ref.TypeParams, err = resolveTypeParams(ctx, ref.TypeParams, actual)
	if err != nil {
		return nil, err
	}

	ref.Methods, err = resolveFields(ctx, ref.Methods, actual)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func resolveTypeDef(ctx Context, ref *types.TypeDefRef, actual types.TypeRefs) (*types.TypeDefRef, error) {
	var err error

	ref.TypeParams, err = resolveTypeParams(ctx, ref.TypeParams, actual)
	if err != nil {
		return nil, err
	}

	ref.Type, err = resolveRef(ctx, ref.Type, actual)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func resolveFuncType(ctx Context, ref *types.FuncTypeRef, actual types.TypeRefs) (*types.FuncTypeRef, error) {
	var err error

	if ref.Params, err = resolveFields(ctx, ref.Params, actual); err != nil {
		return nil, err
	}

	if ref.Results, err = resolveFields(ctx, ref.Results, actual); err != nil {
		return nil, err
	}

	return ref, nil
}

func resolveTypeParam(ctx Context, ref *types.TypeParamRef, actual types.TypeRefs) (types.TypeRef, error) {
	defined, ok := ctx.DefinedByName(ref.Name)
	if !ok {
		return nil, errors.Errorf(ref.Position, "type parameter not found")
	}

	if len(actual)-1 < defined.Order {
		return nil, errors.Errorf(ref.Position, "type parameter not found")
	}

	return actual[defined.Order], nil
}

func resolveTypeParams(ctx Context, params types.TypeRefs, actual types.TypeRefs) (types.TypeRefs, error) {
	for i, param := range params {
		var (
			resolved types.TypeRef
			err      error
		)

		switch p := param.(type) {
		case *types.TypeParamRef:
			resolved, err = resolveTypeParam(ctx, p, actual)
		default:
			resolved, err = resolveRef(ctx, p, actual)
		}
		if err != nil {
			return nil, err
		}
		params[i] = resolved
	}
	return params, nil
}

func resolveFields(ctx Context, ff types.Fields, actual types.TypeRefs) (types.Fields, error) {
	var err error

	for _, field := range ff {
		field.Type, err = resolveRef(ctx, field.Type, actual)
		if err != nil {
			return nil, err
		}
	}

	return ff, nil
}
