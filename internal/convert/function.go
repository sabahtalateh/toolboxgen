package convert

import (
	"fmt"
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	"github.com/sabahtalateh/toolboxgen/internal/mid/position"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) Function(ctx Context, f *ast.FuncDecl) (*types.Function, error) {
	var (
		function *types.Function
		defined  types.TypeParams
		err      error
	)

	function = &types.Function{
		Declared: code.OfNode(f),
		Package:  ctx.Package(),
		FuncName: f.Name.Name,
		Position: ctx.NodePosition(f),
	}

	function.Receiver, defined, err = c.receiver(ctx, f.Recv)
	if err != nil {
		return nil, err
	}

	if function.Receiver == nil {
		function.TypeParams = TypeParams(ctx, f.Type.TypeParams)
		defined = function.TypeParams
	}

	ctx = ctx.WithDefined(defined)

	function.Parameters, err = c.Fields(ctx, f.Type.Params)
	if err != nil {
		return nil, err
	}

	function.Results, err = c.Fields(ctx, f.Type.Results)
	if err != nil {
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
	midType := mid.ParseTypeRef(ctx.Files(), recvField.Type)
	if err = midType.ParseError(); err != nil {
		return nil, nil, err
	}
	defined, err := receiverTypeParams(midType)
	fields, err := c.Field(ctx.WithDefined(defined), recvField)
	if err != nil {
		return nil, nil, err
	}

	return fields[0], defined, err
}

func receiverTypeParams(recv mid.TypeRef) (types.TypeParams, error) {
	var (
		params []mid.TypeRef
		res    types.TypeParams
	)

	switch r := recv.(type) {
	case *mid.Type:
		params = r.TypeParams
	default:
		return nil, errors.Errorf(position.OfTypeRef(recv), "unsupported receiver type")
	}

	for i, param := range params {
		switch p := param.(type) {
		case *mid.Type:
			res = append(res, &types.TypeParam{
				Declared: p.Declared,
				Original: p.TypeName,
				Name:     fmt.Sprintf("T%d", i+1),
				Order:    i,
				Position: p.Position,
			})
		default:
			return nil, errors.Errorf(position.OfTypeRef(p), "unsupported receiver type param type")
		}
	}

	return res, nil
}
