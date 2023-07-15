package tool

// type StructDef struct {
// 	Code     string
// 	Package  string
// 	TypeName string
// 	// Modifiers  []parse.TypeRefModifier
// 	TypeParams []TypeParam
// 	Position   token.Position
// }
//
// func (s StructDef) typDef() {}
//
// type StructRef struct {
// 	Code     string
// 	Package  string
// 	TypeName string
// 	// Mods       []parse.TypeRefModifier
// 	TypeParams struct {
// 		Params    []TypeParam
// 		Effective []TypeRef
// 	}
// 	Position token.Position
// }
//
// // func StructRefFromDef(d *StructDef) *StructRef {
// // 	r := &StructRef{
// // 		Code:     d.Code,
// // 		Package:  d.Package,
// // 		TypeName: d.TypeName,
// // 		Mods:     d.Modifiers,
// // 		TypeParams: struct {
// // 			Params    []TypeParam
// // 			Effective []TypeRef
// // 		}{},
// // 		Position: d.Position,
// // 	}
// //
// // 	for _, param := range d.TypeParams {
// // 		r.TypeParams.Params = append(r.TypeParams.Params, param)
// // 		r.TypeParams.Effective = append(r.TypeParams.Effective, nil)
// // 	}
// //
// // 	return r
// // }
//
// func (s *StructRef) typRef() {}
//
// func (s *StructRef) Equals(t TypeRef) bool {
// 	switch t2 := t.(type) {
// 	case *StructRef:
// 		if !parse.ModifiersEquals(s.Modifiers(), t2.Modifiers()) {
// 			return false
// 		}
//
// 		if s.Package != t2.Package {
// 			return false
// 		}
//
// 		if s.TypeName != t2.TypeName {
// 			return false
// 		}
//
// 		if len(s.TypeParams.Effective) != len(t2.TypeParams.Effective) {
// 			return false
// 		}
//
// 		for i, tp := range s.TypeParams.Effective {
// 			if !tp.Equals(t2.TypeParams.Effective[i]) {
// 				return false
// 			}
// 		}
//
// 		return true
// 	default:
// 		return false
// 	}
// }
//
// func (s *StructRef) Modifiers() []parse.TypeRefModifier {
// 	return s.Mods
// }
//
// func (s *StructRef) NthTypeParam(n int) (*TypeParam, error) {
// 	return nthTypeParam(&s.TypeParams, n)
// }
//
// func (s *StructRef) NumberOfTypeParams() int {
// 	return numberOfParams(&s.TypeParams)
// }
//
// func (s *StructRef) NormalizeTypeParams(normMap map[string]string) {
// 	normalizeTypeParams(&s.TypeParams, normMap)
// }
//
// func (s *StructRef) SetEffectiveParamRecursive(param string, typ TypeRef) {
// 	setEffectiveTypeParamRec(&s.TypeParams, param, typ)
// }
//
// func (s *StructRef) RenameTypeParam(old string, new string) {
// 	renameTypeParam(&s.TypeParams, old, new)
// }
//
// func (s *StructRef) renameTypeParamRecursive(old string, new string) {
// 	renameTypeParamRecursive(&s.TypeParams, old, new)
// }
