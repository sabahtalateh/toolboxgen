package tool

import (
	"go/token"

	"golang.design/x/reflect"
)

type FuncParam struct {
	Name     string
	Type     TypeRef
	Position token.Position
}

type Receiver struct {
	Presented  bool
	Type       TypeRef
	TypeParams []TypeParam
	Position   token.Position
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

// func (f *FuncDef) Clone() *FuncDef {
// 	fd := new(FuncDef)
//
// 	fd.Package = f.Package
// 	fd.FuncName = f.FuncName
//
// 	recv := Receiver{
// 		Presented: f.Receiver.Presented,
// 		Type:      nil,
// 		Position:  f.Receiver.Position,
// 	}
// 	copy(recv.TypeParams, f.Receiver.TypeParams)
// 	fd.Receiver = recv
//
// 	return fd
// }

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

// FuncRefFromDef
// Deep-copying fields to keep original function definition when set effective type parameters
func FuncRefFromDef(pos token.Position, def *FuncDef) *FuncRef {
	fr := &FuncRef{
		Def:      def,
		Package:  def.Package,
		FuncName: def.FuncName,
		TypeParams: struct {
			Params    []TypeParam
			Effective []TypeRef
		}{
			Params:    reflect.DeepCopy[[]TypeParam](def.TypeParams),
			Effective: make([]TypeRef, len(def.TypeParams)),
		},
		Parameters: reflect.DeepCopy[[]FuncParam](def.Parameters),
		Results:    reflect.DeepCopy[[]TypeRef](def.Results),
		Position:   pos,
	}

	return fr
}

func (f *FuncRef) NthTypeParam(n int) (*TypeParam, error) {
	return nthTypeParam(&f.TypeParams, n)
}

func (f *FuncRef) SetEffectiveParam(typeParamName string, effectiveType TypeRef) {
	setEffectiveTypeParam(&f.TypeParams, typeParamName, effectiveType)

	for i := range f.Parameters {
		if _, ok := f.Parameters[i].Type.(*TypeParamRef); ok {
			f.Parameters[i].Type = effectiveType
		}
	}
	for _, par := range f.Parameters {
		if p, ok := par.Type.(ParametrizedRef); ok {
			p.SetEffectiveParam(typeParamName, effectiveType)
		}
	}

	for i := range f.Results {
		if _, ok := f.Results[i].(*TypeParamRef); ok {
			f.Results[i] = effectiveType
		}
	}
	for _, res := range f.Results {
		if r, ok := res.(ParametrizedRef); ok {
			r.SetEffectiveParam(typeParamName, effectiveType)
		}
	}
}

func (f *FuncRef) NumberOfTypeParams() int {
	return numberOfParams(&f.TypeParams)
}

func (f *FuncRef) RenameTypeParam(old string, new string) {
	renameTypeParam(&f.TypeParams, old, new)
}
