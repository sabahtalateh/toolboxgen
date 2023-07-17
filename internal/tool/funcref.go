package tool

// type FuncRef struct {
// 	Code       string
// 	Def        *FuncDef
// 	pakage    string
// 	FuncName   string
// 	TypeParams struct {
// 		Params    []TypeParam
// 		Effective []TypeRef
// 	}
// 	Parameters []FuncParam
// 	Results    []TypeRef
// 	NodePosition   token.NodePosition
// }
//
// // FuncRefFromDef
// // Deep-copying fields to keep original function definition when set effective type parameters
// func FuncRefFromDef(pos token.NodePosition, def *FuncDef) *FuncRef {
// 	fr := &FuncRef{
// 		Def:      def,
// 		pakage:  def.pakage,
// 		FuncName: def.FuncName,
// 		TypeParams: struct {
// 			Params    []TypeParam
// 			Effective []TypeRef
// 		}{
// 			Params:    reflect.DeepCopy[[]TypeParam](def.TypeParams),
// 			Effective: make([]TypeRef, len(def.TypeParams)),
// 		},
// 		Parameters: reflect.DeepCopy[[]FuncParam](def.Parameters),
// 		Results:    reflect.DeepCopy[[]TypeRef](def.Results),
// 		NodePosition:   pos,
// 	}
//
// 	return fr
// }
//
// func (f *FuncRef) NthTypeParam(n int) (*TypeParam, error) {
// 	return nthTypeParam(&f.TypeParams, n)
// }
//
// func (f *FuncRef) NormalizeTypeParams(normMap map[string]string) {
// 	normalizeTypeParams(&f.TypeParams, normMap)
// }
//
// func (f *FuncRef) SetEffectiveParamRecursive(typeParamName string, effectiveType TypeRef) {
// 	setEffectiveTypeParamRec(&f.TypeParams, typeParamName, effectiveType)
//
// 	for i, p := range f.Parameters {
// 		switch pp := p.Type.(type) {
// 		case *TypeParamRef:
// 			if pp.Name != typeParamName {
// 				continue
// 			}
// 			f.Parameters[i].Type = prependModifiers(pp.Modifiers(), effectiveType)
// 		case ParametrizedRef:
// 			pp.SetEffectiveParamRecursive(typeParamName, effectiveType)
// 		}
// 	}
//
// 	for i, r := range f.Results {
// 		switch rr := r.(type) {
// 		case *TypeParamRef:
// 			if rr.Name != typeParamName {
// 				continue
// 			}
// 			f.Results[i] = prependModifiers(rr.Modifiers(), effectiveType)
// 		case ParametrizedRef:
// 			rr.SetEffectiveParamRecursive(typeParamName, effectiveType)
// 		}
// 	}
// }
//
// func (f *FuncRef) NumberOfTypeParams() int {
// 	return numberOfParams(&f.TypeParams)
// }
//
// func (f *FuncRef) RenameTypeParam(old string, new string) {
// 	renameTypeParam(&f.TypeParams, old, new)
// }
//
// func (f *FuncRef) renameTypeParamRecursive(old string, new string) {
// 	// renameTypeParamRecursive(&f.TypeParams, old, new)
// }