package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/maps"
	// "github.com/sabahtalateh/toolboxgen/internal/mid"
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) Function(ctx Context, f *ast.FuncDecl) (*types.Function, error) {
	var (
		function *types.Function
		err      error
	)

	function = &types.Function{
		Declared: code.OfNode(f),
		Package:  ctx.Package,
		FuncName: f.Name.Name,
		Position: ctx.Position(f.Pos()),
	}

	// function.Receiver, err = c.receiver(ctx, f.Recv)
	// if err != nil {
	// 	return nil, err
	// }

	function.TypeParams = TypeParams(ctx, f.Type.TypeParams)
	defined := maps.FromSlice(function.TypeParams, func(v *types.TypeParam) string { return v.Original })
	ctx = ctx.WithDefined(defined)

	function.Parameters, err = c.Fields(ctx.WithPos(f.Type.Params.Pos()), f.Type.Params)
	if err != nil {
		return nil, err
	}

	// function.Results, err = c.Fields(ctx.WithPos(f.Type.Results.Pos()), f.Type.Results)
	// if err != nil {
	// 	return nil, err
	// }

	return function, err
}

func (c *Converter) receiver(ctx Context, r *ast.FieldList) (*types.Field, error) {
	if r == nil || len(r.List) == 0 {
		return nil, nil
	}

	var (
		recv *types.Field
		err  error
	)

	recvField := r.List[0]
	// TODO actual type params
	fields, err := c.Field(ctx, recvField)
	if err != nil {
		return recv, err
	}

	return fields[0], err
}
