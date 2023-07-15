package types

import "go/token"

type (
	TypeParam struct {
		Declared string
		Original string
		Name     string
		Position token.Position
	}

	TypeParams []*TypeParam
)
