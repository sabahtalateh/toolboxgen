package tool

import "go/token"

type FuncParam struct {
	Name     string
	Type     TypeRef
	Position token.Position
}

type Receiver struct {
	WithReceiver bool
	Type         TypeRef
	TypeParams   []TypeParam
	Position     token.Position
}

type FuncDef struct {
	Package    string
	FuncName   string
	Receiver   Receiver
	TypeParams []TypeParam
	Parameters []FuncParam
	Results    []TypeRef
	Position   token.Position
}

type FuncRef struct {
	Def        *FuncDef
	Package    string
	FuncName   string
	TypeParams struct {
		Params    []TypeParam
		Effective []TypeRef
	}
	Parameters []FuncParam
	Results    []TypeRef
	Position   token.Position
}

func FuncRefFromDef(pos token.Position, def *FuncDef) *FuncRef {
	fr := &FuncRef{
		Def:      def,
		Package:  def.Package,
		FuncName: def.FuncName,
		TypeParams: struct {
			Params    []TypeParam
			Effective []TypeRef
		}{},
		Parameters: def.Parameters,
		Results:    def.Results,
		Position:   pos,
	}

	for _, p := range def.TypeParams {
		fr.TypeParams.Params = append(fr.TypeParams.Params, p)
		fr.TypeParams.Effective = append(fr.TypeParams.Effective, nil)
	}

	return fr
}

func (f *FuncRef) NthTypeParam(n int) (*TypeParam, error) {
	return nthTypeParam(&f.TypeParams, n)
}

func (f *FuncRef) SetEffectiveParam(param string, typ TypeRef) {
	setTypeParamRecursive(&f.TypeParams, param, typ)

	for _, par := range f.Parameters {
		if p, ok := par.Type.(ParametrizedRef); ok {
			p.SetEffectiveParam(param, typ)
		}
	}

	for _, ret := range f.Results {
		if p, ok := ret.(ParametrizedRef); ok {
			p.SetEffectiveParam(param, typ)
		}
	}
}

func (f *FuncRef) NumberOfTypeParams() int {
	return numberOfParams(&f.TypeParams)
}

func (f *FuncRef) RenameTypeParam(old string, new string) {
	renameTypeParam(&f.TypeParams, old, new)
}
