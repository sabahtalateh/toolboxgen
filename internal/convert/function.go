package convert

import (
	"fmt"
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) Function(ctx Context, f *ast.FuncDecl) (*types.Function, error) {
	var (
		function *types.Function
		defined  types.TypeParams
		err      error
	)

	function = &types.Function{
		Package:  ctx.Package(),
		FuncName: f.Name.Name,
		Position: ctx.NodePosition(f),
		Code:     code.OfNode(f),
	}

	if function.Receiver, defined, err = c.receiver(ctx, f.Recv); err != nil {
		return nil, err
	}

	if function.Receiver == nil {
		function.TypeParams = TypeParams(ctx, f.Type.TypeParams)
		defined = function.TypeParams
	}

	ctx = ctx.WithDefined(defined)

	if function.Parameters, err = c.Fields(ctx, f.Type.Params); err != nil {
		return nil, err
	}

	if function.Results, err = c.Fields(ctx, f.Type.Results); err != nil {
		return nil, err
	}

	return function, err
}

func (c *Converter) receiver(ctx Context, r *ast.FieldList) (*types.Field, types.TypeParams, error) {
	if r == nil || len(r.List) == 0 {
		return nil, nil, nil
	}

	var err error

	recvField := r.List[0]

	typeExpr := syntax.ParseTypeExpr(ctx.Files(), recvField.Type)
	if err = typeExpr.Error(); err != nil {
		return nil, nil, err
	}

	defined, err := receiverTypeParams(typeExpr)
	if err != nil {
		return nil, nil, err
	}

	fields, err := c.Field(ctx.WithDefined(defined), recvField)
	if err != nil {
		return nil, nil, err
	}

	return fields[0], defined, err
}

func receiverTypeParams(recv syntax.TypeExpr) (types.TypeParams, error) {
	var (
		params []syntax.TypeExpr
		res    types.TypeParams
	)

	switch r := recv.(type) {
	case *syntax.Type:
		params = r.TypeParams
	default:
		return nil, errors.Errorf(recv.Get().Position(), "unsupported receiver type")
	}

	for i, param := range params {
		switch p := param.(type) {
		case *syntax.Type:
			res = append(res, &types.TypeParam{
				Order:    i,
				Name:     fmt.Sprintf("T%d", i+1),
				Position: p.Position,
				Code:     p.Code,
			})
		default:
			return nil, errors.Errorf(p.Get().Position(), "unsupported receiver type param type")
		}
	}

	return res, nil
}
