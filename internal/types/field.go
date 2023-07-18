package types

import "go/token"

type (
	Field struct {
		Declared string
		Name     string
		Type     TypeRef
		Position token.Position
	}

	Fields []*Field
)

func (f *Field) Equal(f2 *Field) bool {
	if f.Name != f2.Name {
		return false
	}

	if !f.Type.Equal(f2.Type) {
		return false
	}

	return true
}

func (f Fields) Equal(f2 Fields) bool {
	if len(f) != len(f2) {
		return false
	}

	for i, field := range f {
		if !field.Equal(f2[i]) {
			return false
		}
	}

	return true
}
