package tool

import "go/token"

type BuiltinRef struct {
	TypeName string
	Pointer  bool
	Position token.Position
}

func (b *BuiltinRef) typ() {
}

func (b *BuiltinRef) Equals(t TypeRef) bool {
	switch t2 := t.(type) {
	case *BuiltinRef:
		if b.TypeName != t2.TypeName {
			return false
		}

		return true
	default:
		return false
	}
}

func (b *BuiltinRef) IsError() bool {
	return b.TypeName == "error"
}
