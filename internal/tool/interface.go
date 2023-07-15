package tool

// type InterfaceDef struct {
// 	Code     string
// 	Package  string
// 	TypeName string
// 	// Modifiers  []parse.TypeRefModifier
// 	TypeParams []TypeParam
// 	Position   token.Position
// }
//
// func (i InterfaceDef) typDef() {}
//
// type InterfaceRef struct {
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
// func InterfaceRefFromDef(d *InterfaceDef) *InterfaceRef {
// 	r := &InterfaceRef{
// 		Code:     d.Code,
// 		Package:  d.Package,
// 		TypeName: d.TypeName,
// 		Mods:     d.Modifiers,
// 		TypeParams: struct {
// 			Params    []TypeParam
// 			Effective []TypeRef
// 		}{},
// 		Position: d.Position,
// 	}
//
// 	for _, param := range d.TypeParams {
// 		r.TypeParams.Params = append(r.TypeParams.Params, param)
// 		r.TypeParams.Effective = append(r.TypeParams.Effective, nil)
// 	}
//
// 	return r
// }
//
// func (i *InterfaceRef) typRef() {}
//
// func (i *InterfaceRef) Equals(t TypeRef) bool {
// 	switch t2 := t.(type) {
// 	case *InterfaceRef:
// 		if !parse.ModifiersEquals(i.Modifiers(), t2.Modifiers()) {
// 			return false
// 		}
//
// 		if i.Package != t2.Package {
// 			return false
// 		}
//
// 		if i.TypeName != t2.TypeName {
// 			return false
// 		}
//
// 		if len(i.TypeParams.Effective) != len(t2.TypeParams.Effective) {
// 			return false
// 		}
//
// 		for idx, tp := range i.TypeParams.Effective {
// 			if !tp.Equals(t2.TypeParams.Effective[idx]) {
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
// func (i *InterfaceRef) Modifiers() []parse.TypeRefModifier {
// 	return i.Mods
// }
//
// func (i *InterfaceRef) NthTypeParam(n int) (*TypeParam, error) {
// 	return nthTypeParam(&i.TypeParams, n)
// }
//
// func (i *InterfaceRef) NumberOfTypeParams() int {
// 	return numberOfParams(&i.TypeParams)
// }
//
// func (i *InterfaceRef) NormalizeTypeParams(normMap map[string]string) {
// 	normalizeTypeParams(&i.TypeParams, normMap)
// }
//
// func (i *InterfaceRef) SetEffectiveParamRecursive(param string, typ TypeRef) {
// 	setEffectiveTypeParamRec(&i.TypeParams, param, typ)
// }
//
// func (i *InterfaceRef) RenameTypeParam(old string, new string) {
// 	renameTypeParam(&i.TypeParams, old, new)
// }
//
// func (i *InterfaceRef) renameTypeParamRecursive(old string, new string) {
// 	renameTypeParamRecursive(&i.TypeParams, old, new)
// }
