package syntax

import "go/token"

type (
	Modifier interface {
		modifier()
	}

	Pointer struct {
		Position token.Position
	}

	Reference struct {
		Position token.Position
	}

	Array struct {
		Sized    bool
		Position token.Position
	}

	Ellipsis struct {
		Position token.Position
	}

	Modifiers []Modifier
)

func (p *Pointer) modifier()   {}
func (p *Reference) modifier() {}
func (a *Array) modifier()     {}
func (a *Ellipsis) modifier()  {}
