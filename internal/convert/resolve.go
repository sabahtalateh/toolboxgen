package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

// resolveRef resolves type parameters into actual types
func resolveRef(ctx Context, ref types.TypeRef, actual types.TypeRefs) (types.TypeRef, error) {
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
		return r, nil
	case *types.MapRef:
		return resolveMap(ctx, r, actual)
	case *types.ChanRef:
		return resolveChan(ctx, r, actual)
	case *types.FuncTypeRef:
		return resolveFuncType(ctx, r, actual)
	case *types.StructTypeRef:
		return resolveStructType(ctx, r, actual)
	case *types.InterfaceTypeRef:
		return resolveInterfaceType(ctx, r, actual)
	case *types.TypeParamRef:
		return resolveTypeParam(ctx, r, actual)
	default:
		return nil, errors.Errorf(r.Get().Position(), "unknown type %T", r)
	}
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

	ref.Fields, err = resolveFields(ctx, ref.Fields, actual)
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

func resolveMap(ctx Context, ref *types.MapRef, actual types.TypeRefs) (*types.MapRef, error) {
	var err error

	ref.Key, err = resolveRef(ctx, ref.Key, actual)
	if err != nil {
		return nil, err
	}

	ref.Value, err = resolveRef(ctx, ref.Value, actual)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func resolveChan(ctx Context, ref *types.ChanRef, actual types.TypeRefs) (*types.ChanRef, error) {
	var err error

	ref.Value, err = resolveRef(ctx, ref.Value, actual)
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

func resolveStructType(ctx Context, ref *types.StructTypeRef, actual types.TypeRefs) (*types.StructTypeRef, error) {
	var err error

	if ref.Fields, err = resolveFields(ctx, ref.Fields, actual); err != nil {
		return nil, err
	}

	return ref, nil
}

func resolveInterfaceType(ctx Context, ref *types.InterfaceTypeRef, actual types.TypeRefs) (*types.InterfaceTypeRef, error) {
	var err error

	if ref.Fields, err = resolveFields(ctx, ref.Fields, actual); err != nil {
		return nil, err
	}

	return ref, nil
}

func resolveTypeParam(ctx Context, ref *types.TypeParamRef, actual types.TypeRefs) (types.TypeRef, error) {
	actual = actual.Clone()

	defined, ok := ctx.DefinedByName(ref.Name)
	if !ok {
		return nil, errors.Errorf(ref.Position, "type parameter not found")
	}

	if len(actual)-1 < defined.Order {
		return nil, errors.Errorf(ref.Position, "type parameter not found")
	}

	a := actual[defined.Order]
	a.Set().Modifiers(append(ref.Modifiers, a.Get().Modifiers()...))

	return a, nil
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
