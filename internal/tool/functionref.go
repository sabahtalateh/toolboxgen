package tool

import "go/token"

type FunctionRef struct {
	Package    string
	FuncName   string
	TypeParams struct {
		Params    []TypeParam
		Effective []TypeRef
	}
	Parameters []TypeRef
	Results    []TypeRef
	Position   token.Position
}

func (f *FunctionRef) NthTypeParam(n int) (*TypeParam, error) {
	return nthTypeParam(&f.TypeParams, n)
}

func (f *FunctionRef) SetEffectiveParam(param string, typ TypeRef) {
	setTypeParamRecursive(&f.TypeParams, param, typ)

	for _, par := range f.Parameters {
		if p, ok := par.(Parametrized); ok {
			p.SetEffectiveParam(param, typ)
		}
	}

	for _, ret := range f.Results {
		if p, ok := ret.(Parametrized); ok {
			p.SetEffectiveParam(param, typ)
		}
	}
}

func (f *FunctionRef) NumberOfParams() int {
	return numberOfParams(&f.TypeParams)
}

func (f *FunctionRef) RenameTypeParam(old string, new string) {
	renameTypeParam(&f.TypeParams, old, new)
}
