package tool

import (
	"errors"
	"go/token"
)

type TypeParam struct {
	Name     string
	Position token.Position
}

// ParametrizedRef struct reference, function reference or interface or
type ParametrizedRef interface {
	NumberOfTypeParams() int
	NthTypeParam(n int) (*TypeParam, error)
	RenameTypeParam(old string, new string)
	SetEffectiveParam(paramName string, typ TypeRef)
}

func setEffectiveTypeParam(typeParams *struct {
	Params    []TypeParam
	Effective []TypeRef
}, typeParamName string, effective TypeRef) {
	idx := -1
	for i, tp := range typeParams.Params {
		if tp.Name == typeParamName {
			idx = i
			break
		}
	}

	if idx != -1 {
		// if typeParams.Effective[idx] != nil {
		// 	return
		// }
		typeParams.Effective[idx] = effective
	}

	// for _, t := range typeParams.Effective {
	// 	tt, ok := t.(ParametrizedRef)
	// 	if ok {
	// 		tt.SetEffectiveParam(typeParamName, effective)
	// 	}
	// }
}

func nthTypeParam(typeParams *struct {
	Params    []TypeParam
	Effective []TypeRef
}, n int) (*TypeParam, error) {
	if n > len(typeParams.Params)-1 {
		return nil, errors.New("type param index out of bounds")
	}
	return &typeParams.Params[n], nil
}

func numberOfParams(typeParams *struct {
	Params    []TypeParam
	Effective []TypeRef
}) int {
	return len(typeParams.Params)
}

func renameTypeParam(typeParams *struct {
	Params    []TypeParam
	Effective []TypeRef
}, old string, new string) {
	for i, p := range typeParams.Params {
		if p.Name == old {
			typeParams.Params[i] = TypeParam{
				Name:     new,
				Position: p.Position,
			}
		}
	}
}
