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
			newName := fmt.Sprintf("T%d", i+1)
			res = append(res, &types.TypeParam{
				Declared: code.OfNode(name),
				Original: name.Name,
				Name:     newName,
				Order:    i,
				Position: ctx.NodePosition(name),
			})
			i += 1
		}
	}

	return res
}
