package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

// resolveType resolves type parameters into actual types
func resolveType(ctx Context, t types.TypeRef, actual types.TypeRefs) (types.TypeRef, error) {
	switch r := t.(type) {
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

func resolveStruct(ctx Context, t *types.StructRef, actual types.TypeRefs) (*types.StructRef, error) {
	var err error

	t.TypeParams, err = resolveTypeParams(ctx, t.TypeParams, actual)
	if err != nil {
		return nil, err
	}

	t.Fields, err = resolveFields(ctx, t.Fields, actual)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func resolveInterface(ctx Context, t *types.InterfaceRef, actual types.TypeRefs) (*types.InterfaceRef, error) {
	var err error

	t.TypeParams, err = resolveTypeParams(ctx, t.TypeParams, actual)
	if err != nil {
		return nil, err
	}

	t.Fields, err = resolveFields(ctx, t.Fields, actual)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func resolveTypeDef(ctx Context, t *types.TypeDefRef, actual types.TypeRefs) (*types.TypeDefRef, error) {
	var err error

	t.TypeParams, err = resolveTypeParams(ctx, t.TypeParams, actual)
	if err != nil {
		return nil, err
	}

	t.Type, err = resolveType(ctx, t.Type, actual)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func resolveMap(ctx Context, t *types.MapRef, actual types.TypeRefs) (*types.MapRef, error) {
	var err error

	t.Key, err = resolveType(ctx, t.Key, actual)
	if err != nil {
		return nil, err
	}

	t.Value, err = resolveType(ctx, t.Value, actual)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func resolveChan(ctx Context, t *types.ChanRef, actual types.TypeRefs) (*types.ChanRef, error) {
	var err error

	t.Value, err = resolveType(ctx, t.Value, actual)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func resolveFuncType(ctx Context, t *types.FuncTypeRef, actual types.TypeRefs) (*types.FuncTypeRef, error) {
	var err error

	if t.Params, err = resolveFields(ctx, t.Params, actual); err != nil {
		return nil, err
	}

	if t.Results, err = resolveFields(ctx, t.Results, actual); err != nil {
		return nil, err
	}

	return t, nil
}

func resolveStructType(ctx Context, t *types.StructTypeRef, actual types.TypeRefs) (*types.StructTypeRef, error) {
	var err error

	if t.Fields, err = resolveFields(ctx, t.Fields, actual); err != nil {
		return nil, err
	}

	return t, nil
}

func resolveInterfaceType(ctx Context, t *types.InterfaceTypeRef, actual types.TypeRefs) (*types.InterfaceTypeRef, error) {
	var err error

	if t.Fields, err = resolveFields(ctx, t.Fields, actual); err != nil {
		return nil, err
	}

	return t, nil
}

func resolveTypeParam(ctx Context, t *types.TypeParamRef, actual types.TypeRefs) (types.TypeRef, error) {
	actual = actual.Clone()

	defined, ok := ctx.DefinedByName(t.Name)
	if !ok {
		return nil, errors.Errorf(t.Position, "type parameter not found")
	}

	if len(actual)-1 < defined.Order {
		return nil, errors.Errorf(t.Position, "type parameter not found")
	}

	a := actual[defined.Order]
	a.Set().Modifiers(append(t.Modifiers, a.Get().Modifiers()...))

	return a, nil
}

func resolveTypeParams(ctx Context, tt types.TypeRefs, actual types.TypeRefs) (types.TypeRefs, error) {
	for i, param := range tt {
		var (
			resolved types.TypeRef
			err      error
		)

		switch p := param.(type) {
		case *types.TypeParamRef:
			resolved, err = resolveTypeParam(ctx, p, actual)
		default:
			resolved, err = resolveType(ctx, p, actual)
		}
		if err != nil {
			return nil, err
		}
		tt[i] = resolved
	}
	return tt, nil
}

func resolveFields(ctx Context, ff types.Fields, actual types.TypeRefs) (types.Fields, error) {
	var err error

	for _, field := range ff {
		field.Type, err = resolveType(ctx, field.Type, actual)
		if err != nil {
			return nil, err
		}
	}

	return ff, nil
}
