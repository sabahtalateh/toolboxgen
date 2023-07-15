package types

import "go/token"

type (
	Modifier interface {
		modifier()
		Equal(m Modifier) bool
	}

	Modifiers []Modifier

	Pointer struct {
		Position token.Position
	}

	Array struct {
		Sized    bool
		Position token.Position
	}

	Ellipsis struct {
		Position token.Position
	}
)

func (m *Pointer) modifier()  {}
func (m *Array) modifier()    {}
func (m *Ellipsis) modifier() {}

func (m *Pointer) Equal(m2 Modifier) bool {
	switch m2.(type) {
	case *Pointer:
		return true
	default:
		return false
	}
}

func (m *Pointer) String() string {
	return "*"
}

func (m *Array) Equal(m2 Modifier) bool {
	switch m2.(type) {
	case *Array:
		// don't care of array size
		return true
	default:
		return false
	}
}

func (m *Array) String() string {
	return "[]"
}

func (m *Ellipsis) Equal(m2 Modifier) bool {
	switch m2.(type) {
	case *Ellipsis:
		return true
	default:
		return false
	}
}

func (m *Ellipsis) String() string {
	return "..."
}
