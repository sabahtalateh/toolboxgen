package convert

import (
	"fmt"
	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"go/ast"
)

func TypeParams(ctx Context, ff *ast.FieldList) types.TypeParams {
	if ff == nil {
		return nil
	}

	var (
		res types.TypeParams
		i   = 1
	)

	for _, field := range ff.List {
		for _, name := range field.Names {
			newName := fmt.Sprintf("T%d", i)
			res = append(res, &types.TypeParam{
				Declared: code.OfNode(name),
				Original: name.Name,
				Name:     newName,
				Position: ctx.Position(name.Pos()),
			})
			i += 1
		}
	}

	return res
}

func TypeParamsRename(pp types.TypeParams) map[string]string {
	res := make(map[string]string)

	for _, param := range pp {
		res[param.Original] = param.Name
	}

	return res
}

func TypeRefRename(r types.TypeRef) map[string]string {
	res := make(map[string]string)

	var typeParams []types.TypeRef
	switch t := r.(type) {
	case *types.StructRef:
		typeParams = t.TypeParams
	case *types.InterfaceRef:
		typeParams = t.TypeParams
	case *types.TypeDefRef:
		typeParams = t.TypeParams
	}
	for _, param := range typeParams {
		switch t := param.(type) {
		case *types.TypeParamRef:
			res[t.Original] = t.Name
		}
	}

	return res
}
