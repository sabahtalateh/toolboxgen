package tool

import (
	"errors"
	"go/token"
)

type TypeParam struct {
	Name     string
	Position token.Position
}

type Parametrized interface {
	NumberOfParams() int
	NthTypeParam(n int) (*TypeParam, error)
	RenameTypeParam(old string, new string)
	SetEffectiveParam(paramName string, typ TypeRef)
}

func setTypeParamRecursive(typeParams *struct {
	Params    []TypeParam
	Effective []TypeRef
}, param string, effective TypeRef) {
	idx := -1
	for i, tp := range typeParams.Params {
		if tp.Name == param {
			idx = i
			break
		}
	}
	if idx != -1 {
		typeParams.Effective[idx] = effective
	}

	if typeParams.Params == nil {
		return
	}

	for _, t := range typeParams.Effective {
		tt, ok := t.(Parametrized)
		if ok {
			tt.SetEffectiveParam(param, effective)
		}
	}
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
