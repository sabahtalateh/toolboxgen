package types

import "go/token"

type (
	TypeParam struct {
		Declared string
		Original string
		Name     string
		Order    int
		Position token.Position
	}

	TypeParams []*TypeParam
)

func (t TypeParams) Equal(t2 TypeParams) bool {
	if len(t) != len(t2) {
		return false
	}

	for i, param := range t {
		if !param.Equal(t2[i]) {
			return false
		}
	}

	return true
}

func (t *TypeParam) Equal(t2 *TypeParam) bool {
	return t.Name == t2.Name
}
