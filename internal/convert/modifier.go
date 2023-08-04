package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func Modifiers(modifiers mid.Modifiers) types.Modifiers {
	res := make([]types.Modifier, 0, len(modifiers))

	for _, m := range modifiers {
		switch modifier := m.(type) {
		case *mid.Pointer:
			res = append(res, &types.Pointer{Position: modifier.Position})
		case *mid.Array:
			res = append(res, &types.Array{Sized: modifier.Sized, Position: modifier.Position})
		case *mid.Ellipsis2:
			res = append(res, &types.Ellipsis{Position: modifier.Position})
		}
	}

	return res
}
