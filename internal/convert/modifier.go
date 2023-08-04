package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func Modifiers(modifiers syntax.Modifiers) types.Modifiers {
	res := make([]types.Modifier, 0, len(modifiers))

	for _, m := range modifiers {
		switch modifier := m.(type) {
		case *syntax.Pointer:
			res = append(res, &types.Pointer{Position: modifier.Position})
		case *syntax.Array:
			res = append(res, &types.Array{Sized: modifier.Sized, Position: modifier.Position})
		case *syntax.Ellipsis2:
			res = append(res, &types.Ellipsis{Position: modifier.Position})
		}
	}

	return res
}
