package convert

import (
	"fmt"
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

// Field returns at least 1 element
func (c *Converter) Field(ctx Context, field *ast.Field) (types.Fields, error) {
	tr, err := c.TypeRef(ctx.WithPos(field.Type.Pos()), field.Type)
	if err != nil {
		return nil, err
	}

	if len(field.Names) == 0 {
		return []*types.Field{{
			Type:     tr,
			Position: ctx.NodePosition(field),
			Code:     code.OfNode(field.Type)},
		}, nil
	}

	var res types.Fields
	for _, name := range field.Names {
		res = append(res, &types.Field{
			Name:     name.Name,
			Type:     tr,
			Position: ctx.NodePosition(name),
			Code:     fmt.Sprintf("%s %s", name.Name, code.OfNode(field.Type)),
		})
	}

	return res, nil
}

func (c *Converter) Fields(ctx Context, ff *ast.FieldList) (types.Fields, error) {
	if ff == nil {
		return nil, nil
	}

	var (
		res types.Fields
	)

	for _, field := range ff.List {
		f, err := c.Field(ctx.WithPos(field.Pos()), field)
		if err != nil {
			return nil, err
		}
		res = append(res, f...)
	}

	return res, nil
}
