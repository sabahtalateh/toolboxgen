package syntax

type TypeRefModifier interface {
	typeRefModifier()
	Equal(m2 TypeRefModifier) bool
}

type Pointer struct{}

func (p *Pointer) Equal(m2 TypeRefModifier) bool {
	switch m2.(type) {
	case *Pointer:
		return true
	default:
		return false
	}
}

func (p *Pointer) typeRefModifier() {}

type Array struct {
	Sized bool
}

func (p *Array) Equal(m2 TypeRefModifier) bool {
	switch ar := m2.(type) {
	case *Array:
		// sized arrays not supported
		if !p.Sized || !ar.Sized {
			return false
		}
		return true
	default:
		return false
	}
}

func (a *Array) typeRefModifier() {}

func ModifiersEquals(mm1 []TypeRefModifier, mm2 []TypeRefModifier) bool {
	if len(mm1) != len(mm2) {
		return false
	}

	for i := 0; i < len(mm1); i++ {
		if !mm1[i].Equal(mm2[i]) {
			return false
		}
	}

	return true
}
