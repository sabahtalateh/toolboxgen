package tool

import "go/token"

type StructRef struct {
	Package    string
	TypeName   string
	Pointer    bool
	TypeParams struct {
		Params    []TypeParam
		Effective []TypeRef
	}
	Position token.Position
}

func (s *StructRef) typ() {
}

func (s *StructRef) Equals(t TypeRef) bool {
	switch t2 := t.(type) {
	case *StructRef:
		if s.Pointer != t2.Pointer {
			return false
		}

		if s.Package != t2.Package {
			return false
		}

		if s.TypeName != t2.TypeName {
			return false
		}

		if len(s.TypeParams.Effective) != len(t2.TypeParams.Effective) {
			return false
		}

		for i, tp := range s.TypeParams.Effective {
			if !tp.Equals(t2.TypeParams.Effective[i]) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func (s *StructRef) NthTypeParam(n int) (*TypeParam, error) {
	return nthTypeParam(&s.TypeParams, n)
}

func (s *StructRef) NumberOfParams() int {
	return numberOfParams(&s.TypeParams)
}

func (s *StructRef) SetEffectiveParam(param string, typ TypeRef) {
	setTypeParamRecursive(&s.TypeParams, param, typ)
}

func (s *StructRef) RenameTypeParam(old string, new string) {
	renameTypeParam(&s.TypeParams, old, new)
}
