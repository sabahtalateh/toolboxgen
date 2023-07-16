package convert

import (
	"fmt"
	"go/ast"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

// Field returns at least 1 element
func (c *Converter) Field(ctx Context, field *ast.Field) ([]*types.Field, error) {
	tr, err := c.TypeRef(ctx.WithPos(field.Type.Pos()), field.Type)
	if err != nil {
		return nil, err
	}

	if len(field.Names) == 0 {
		return []*types.Field{{Declared: code.OfNode(field.Type), Type: tr, Position: ctx.NodePosition(field)}}, nil
	}

	var res []*types.Field
	for _, name := range field.Names {
		res = append(res, &types.Field{
			Declared: fmt.Sprintf("%s %s", name.Name, code.OfNode(field.Type)),
			Name:     name.Name,
			Type:     tr,
			Position: ctx.NodePosition(name),
		})
	}

	return res, nil
}

func (c *Converter) Fields(ctx Context, ff *ast.FieldList) ([]*types.Field, error) {
	if ff == nil {
		return nil, nil
	}

	var (
		res []*types.Field
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
