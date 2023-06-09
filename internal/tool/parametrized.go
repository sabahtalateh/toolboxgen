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
	RenameTypeParamRecursive(old string, new string)
	SetEffectiveParamRecursive(paramName string, typ TypeRef)
}

func setEffectiveTypeParamRec(typeParams *struct {
	Params    []TypeParam
	Effective []TypeRef
}, typeParamName string, effective TypeRef) {
	typeParamIdx := -1
	for i, tp := range typeParams.Params {
		if tp.Name == typeParamName {
			typeParamIdx = i
			break
		}
	}

	if typeParamIdx != -1 {
		// if typeParams.Effective[typeParamIdx] != nil {
		// 	return
		// }
		if typeParams.Effective[typeParamIdx] != nil {
			effective = prependModifiers(typeParams.Effective[typeParamIdx].Modifiers(), effective)
		}
		typeParams.Effective[typeParamIdx] = effective
	}

	for i, t := range typeParams.Effective {
		if typeParamIdx == i {
			continue
		}
		tt, ok := t.(ParametrizedRef)
		if ok {
			tt.SetEffectiveParamRecursive(typeParamName, effective)
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

func renameTypeParamRecursive(typeParams *struct {
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

	for _, eff := range typeParams.Effective {
		switch t := eff.(type) {
		case ParametrizedRef:
			t.RenameTypeParamRecursive(old, new)
		case *TypeParamRef:
			if t.Name == old {
				t.Name = new
			}
		}
	}
}
