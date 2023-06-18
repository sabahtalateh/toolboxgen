package tool

import "go/token"

type InterfaceRef struct {
	Package    string
	TypeName   string
	Pointer    bool
	TypeParams struct {
		Params    []TypeParam
		Effective []TypeRef
	}
	Position token.Position
}

func (i *InterfaceRef) typ() {
}

func (i *InterfaceRef) Equals(t TypeRef) bool {
	switch t2 := t.(type) {
	case *InterfaceRef:
		if i.Pointer != t2.Pointer {
			return false
		}

		if i.Package != t2.Package {
			return false
		}

		if i.TypeName != t2.TypeName {
			return false
		}

		if len(i.TypeParams.Effective) != len(t2.TypeParams.Effective) {
			return false
		}

		for idx, tp := range i.TypeParams.Effective {
			if !tp.Equals(t2.TypeParams.Effective[idx]) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func (i *InterfaceRef) NthTypeParam(n int) (*TypeParam, error) {
	return nthTypeParam(&i.TypeParams, n)
}

func (i *InterfaceRef) NumberOfParams() int {
	return numberOfParams(&i.TypeParams)
}

func (i *InterfaceRef) SetEffectiveParam(param string, typ TypeRef) {
	setTypeParamRecursive(&i.TypeParams, param, typ)
}

func (i *InterfaceRef) RenameTypeParam(old string, new string) {
	renameTypeParam(&i.TypeParams, old, new)
}
