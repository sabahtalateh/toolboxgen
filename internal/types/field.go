package types

import "go/token"

type (
	Field struct {
		Declared string
		Name     string
		Type     TypeRef
		Position token.Position
	}
)
