package inspect

import (
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func Modifier(m types.Modifier) string {
	switch m.(type) {
	case *types.Pointer:
		return "*"
	case *types.Array:
		return "[]"
	case *types.Ellipsis:
		return "..."
	default:
		panic("unknown modifier")
	}
}

func Modifiers(m types.Modifiers) string {
	return strings.Join(slices.Map(m, func(el types.Modifier) string { return Modifier(el) }), "")
}
