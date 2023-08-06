package convert

import (
	"fmt"
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func TypeParams(ctx Context, ff *ast.FieldList) types.TypeParams {
	if ff == nil {
		return nil
	}

	var (
		res types.TypeParams
		i   = 0
	)

	for _, field := range ff.List {
		for _, name := range field.Names {
			res = append(res, &types.TypeParam{
				Order:    i,
				Name:     name.Name,
				Position: ctx.NodePosition(name),
				Code:     fmt.Sprintf("%s %s", code.OfNode(name), code.OfNode(field.Type)),
			})
			i += 1
		}
	}

	return res
}

func InitTypeParams(params types.TypeParams) types.TypeExprs {
	var res types.TypeExprs

	for _, param := range params {
		res = append(res, &types.TypeParamExpr{
			Order:      param.Order,
			Name:       param.Name,
			Position:   param.Position,
			Definition: param,
			Code:       param.Code,
		})
	}

	return res
}
