package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

// resolveType resolves type parameters into actual types
func resolveType(ctx Context, t types.TypeExpr, actual types.TypeExprs) (types.TypeExpr, error) {
	switch r := t.(type) {
	case *types.BuiltinExpr:
		return r, nil
	case *types.StructExpr:
		return resolveStruct(ctx, r, actual)
	case *types.InterfaceExpr:
		return resolveInterface(ctx, r, actual)
	case *types.TypeDefExpr:
		return resolveTypeDef(ctx, r, actual)
	case *types.TypeAliasExpr:
		return r, nil
	case *types.MapExpr:
		return resolveMap(ctx, r, actual)
	case *types.ChanExpr:
		return resolveChan(ctx, r, actual)
	case *types.FuncTypeExpr:
		return resolveFuncType(ctx, r, actual)
	case *types.StructTypeExpr:
		return resolveStructType(ctx, r, actual)
	case *types.InterfaceTypeExpr:
		return resolveInterfaceType(ctx, r, actual)
	case *types.TypeParamExpr:
		return resolveTypeParam(ctx, r, actual)
	default:
		return nil, errors.Errorf(r.Get().Position(), "unknown type %T", r)
	}
}

func resolveStruct(ctx Context, t *types.StructExpr, actual types.TypeExprs) (*types.StructExpr, error) {
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

func resolveInterface(ctx Context, t *types.InterfaceExpr, actual types.TypeExprs) (*types.InterfaceExpr, error) {
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

func resolveTypeDef(ctx Context, t *types.TypeDefExpr, actual types.TypeExprs) (*types.TypeDefExpr, error) {
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

func resolveMap(ctx Context, t *types.MapExpr, actual types.TypeExprs) (*types.MapExpr, error) {
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

func resolveChan(ctx Context, t *types.ChanExpr, actual types.TypeExprs) (*types.ChanExpr, error) {
	var err error

	t.Value, err = resolveType(ctx, t.Value, actual)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func resolveFuncType(ctx Context, t *types.FuncTypeExpr, actual types.TypeExprs) (*types.FuncTypeExpr, error) {
	var err error

	if t.Params, err = resolveFields(ctx, t.Params, actual); err != nil {
		return nil, err
	}

	if t.Results, err = resolveFields(ctx, t.Results, actual); err != nil {
		return nil, err
	}

	return t, nil
}

func resolveStructType(ctx Context, t *types.StructTypeExpr, actual types.TypeExprs) (*types.StructTypeExpr, error) {
	var err error

	if t.Fields, err = resolveFields(ctx, t.Fields, actual); err != nil {
		return nil, err
	}

	return t, nil
}

func resolveInterfaceType(ctx Context, t *types.InterfaceTypeExpr, actual types.TypeExprs) (*types.InterfaceTypeExpr, error) {
	var err error

	if t.Fields, err = resolveFields(ctx, t.Fields, actual); err != nil {
		return nil, err
	}

	return t, nil
}

func resolveTypeParam(ctx Context, t *types.TypeParamExpr, actual types.TypeExprs) (types.TypeExpr, error) {
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

func resolveTypeParams(ctx Context, tt types.TypeExprs, actual types.TypeExprs) (types.TypeExprs, error) {
	for i, param := range tt {
		var (
			resolved types.TypeExpr
			err      error
		)

		switch p := param.(type) {
		case *types.TypeParamExpr:
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

func resolveFields(ctx Context, ff types.Fields, actual types.TypeExprs) (types.Fields, error) {
	var err error

	for _, field := range ff {
		field.Type, err = resolveType(ctx, field.Type, actual)
		if err != nil {
			return nil, err
		}
	}

	return ff, nil
}
