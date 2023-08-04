package syntax

import "go/token"

type (
	Modifier interface {
		modifier()
	}

	Pointer struct {
		Position token.Position
	}

	Array struct {
		Sized    bool
		Position token.Position
	}

	Ellipsis2 struct {
		Position token.Position
	}

	Modifiers []Modifier
)

func (p *Pointer) modifier()   {}
func (a *Array) modifier()     {}
func (a *Ellipsis2) modifier() {}
