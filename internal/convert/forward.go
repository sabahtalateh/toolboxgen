package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

// forwardTypeArgs forwards type arguments into types
func forwardTypeArgs(ctx Context, t types.TypeExpr, args types.TypeExprs) (types.TypeExpr, error) {
	switch r := t.(type) {
	case *types.BuiltinExpr:
		return r, nil
	case *types.StructExpr:
		return forwardToStruct(ctx, r, args)
	case *types.InterfaceExpr:
		return forwardToInterface(ctx, r, args)
	case *types.TypeDefExpr:
		return forwardToTypeDef(ctx, r, args)
	case *types.TypeAliasExpr:
		return r, nil
	case *types.MapExpr:
		return forwardToMap(ctx, r, args)
	case *types.ChanExpr:
		return forwardToChan(ctx, r, args)
	case *types.FuncTypeExpr:
		return forwardToFuncType(ctx, r, args)
	case *types.StructTypeExpr:
		return forwardToStructType(ctx, r, args)
	case *types.InterfaceTypeExpr:
		return forwardToInterfaceType(ctx, r, args)
	case *types.TypeArgExpr:
		return forwardToTypeArg(ctx, r, args)
	default:
		return nil, errors.Errorf(r.Get().Position(), "unknown type %T", r)
	}
}

func forwardToStruct(ctx Context, t *types.StructExpr, args types.TypeExprs) (*types.StructExpr, error) {
	var err error

	t.TypeArgs, err = forwardToTypeArgs(ctx, t.TypeArgs, args)
	if err != nil {
		return nil, err
	}

	t.Fields, err = forwardToFields(ctx, t.Fields, args)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToInterface(ctx Context, t *types.InterfaceExpr, args types.TypeExprs) (*types.InterfaceExpr, error) {
	var err error

	t.TypeArgs, err = forwardToTypeArgs(ctx, t.TypeArgs, args)
	if err != nil {
		return nil, err
	}

	t.Fields, err = forwardToFields(ctx, t.Fields, args)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToTypeDef(ctx Context, t *types.TypeDefExpr, args types.TypeExprs) (*types.TypeDefExpr, error) {
	var err error

	t.TypeArgs, err = forwardToTypeArgs(ctx, t.TypeArgs, args)
	if err != nil {
		return nil, err
	}

	t.Type, err = forwardTypeArgs(ctx, t.Type, args)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToMap(ctx Context, t *types.MapExpr, args types.TypeExprs) (*types.MapExpr, error) {
	var err error

	t.Key, err = forwardTypeArgs(ctx, t.Key, args)
	if err != nil {
		return nil, err
	}

	t.Value, err = forwardTypeArgs(ctx, t.Value, args)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToChan(ctx Context, t *types.ChanExpr, args types.TypeExprs) (*types.ChanExpr, error) {
	var err error

	t.Value, err = forwardTypeArgs(ctx, t.Value, args)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToFuncType(ctx Context, t *types.FuncTypeExpr, args types.TypeExprs) (*types.FuncTypeExpr, error) {
	var err error

	if t.Params, err = forwardToFields(ctx, t.Params, args); err != nil {
		return nil, err
	}

	if t.Results, err = forwardToFields(ctx, t.Results, args); err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToStructType(ctx Context, t *types.StructTypeExpr, args types.TypeExprs) (*types.StructTypeExpr, error) {
	var err error

	if t.Fields, err = forwardToFields(ctx, t.Fields, args); err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToInterfaceType(ctx Context, t *types.InterfaceTypeExpr, args types.TypeExprs) (*types.InterfaceTypeExpr, error) {
	var err error

	if t.Fields, err = forwardToFields(ctx, t.Fields, args); err != nil {
		return nil, err
	}

	return t, nil
}

func forwardToTypeArg(ctx Context, t *types.TypeArgExpr, args types.TypeExprs) (types.TypeExpr, error) {
	args = args.Clone()

	defined, ok := ctx.DefinedByName(t.Name)
	if !ok {
		return nil, errors.Errorf(t.Position, "type parameter not found")
	}

	if len(args)-1 < defined.Order {
		return nil, errors.Errorf(t.Position, "type parameter not found")
	}

	a := args[defined.Order]
	a.Set().Modifiers(append(t.Modifiers, a.Get().Modifiers()...))

	return a, nil
}

func forwardToTypeArgs(ctx Context, tt types.TypeExprs, args types.TypeExprs) (types.TypeExprs, error) {
	for i, param := range tt {
		var (
			forwarded types.TypeExpr
			err       error
		)

		switch p := param.(type) {
		case *types.TypeArgExpr:
			forwarded, err = forwardToTypeArg(ctx, p, args)
		default:
			forwarded, err = forwardTypeArgs(ctx, p, args)
		}
		if err != nil {
			return nil, err
		}
		tt[i] = forwarded
	}
	return tt, nil
}

func forwardToFields(ctx Context, ff types.Fields, args types.TypeExprs) (types.Fields, error) {
	var err error

	for _, field := range ff {
		field.Type, err = forwardTypeArgs(ctx, field.Type, args)
		if err != nil {
			return nil, err
		}
	}

	return ff, nil
}
