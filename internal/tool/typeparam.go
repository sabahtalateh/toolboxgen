package tool

import "go/token"

type TypeParamRef struct {
	Name     string
	Pointer  bool
	Position token.Position
}

func (s *TypeParamRef) typ() {
}

func (s *TypeParamRef) Equals(t TypeRef) bool {
	switch t2 := t.(type) {
	case *TypeParamRef:
		if s.Pointer != t2.Pointer {
			return false
		}

		if s.Name != t2.Name {
			return false
		}

		return true
	default:
		return false
	}
}
